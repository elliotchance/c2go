// This file contains functions for transpiling function calls (invocations).

package transpiler

import (
	"fmt"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"

	goast "go/ast"
	"go/token"
)

func getName(firstChild ast.Node) string {
	switch fc := firstChild.(type) {
	case *ast.DeclRefExpr:
		return fc.Name

	case *ast.MemberExpr:
		if types.IsFunction(fc.Type) {
			if decl, ok := fc.Children()[0].(*ast.DeclRefExpr); ok {
				return decl.Name + "." + fc.Name
			}
		}
		return fc.Name

	case *ast.ParenExpr:
		return getName(fc.Children()[0])

	case *ast.UnaryOperator:
		return getName(fc.Children()[0])

	case *ast.ImplicitCastExpr:
		return getName(fc.Children()[0])

	case *ast.CStyleCastExpr:
		return getName(fc.Children()[0])

	default:
		panic(fmt.Sprintf("cannot CallExpr on: %#v", fc))
	}
}

func getNameOfFunctionFromCallExpr(n *ast.CallExpr) (string, error) {
	// The first child will always contain the name of the function being
	// called.
	firstChild, ok := n.Children()[0].(*ast.ImplicitCastExpr)
	if !ok {
		err := fmt.Errorf("unable to use CallExpr: %#v", n.Children()[0])
		return "", err
	}

	return getName(firstChild.Children()[0]), nil
}

// transpileCallExpr transpiles expressions that calls a function, for example:
//
//     foo("bar")
//
// It returns three arguments; the Go AST expression, the C type (that is
// returned by the function) and any error. If there is an error returned you
// can assume the first two arguments will not contain any useful information.
func transpileCallExpr(n *ast.CallExpr, p *program.Program) (
	_ *goast.CallExpr, resultType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {

	functionName, err := getNameOfFunctionFromCallExpr(n)
	if err != nil {
		return nil, "", nil, nil, err
	}
	functionName = util.ConvertFunctionNameFromCtoGo(functionName)

	if functionName == "__builtin_va_start" ||
		functionName == "__builtin_va_end" {
		// ignore function __builtin_va_start, __builtin_va_end
		// see "Variadic functions"
		return nil, "", nil, nil, nil
	}

	// function stdlib.c
	if functionName == "calloc" && len(n.Children()) == 3 {
		var allocType string
		size, _, preStmts, postStmts, err := transpileToExpr(n.Children()[1], p, false)
		if err != nil {
			return nil, "", nil, nil, err
		}
		if v, ok := n.Children()[2].(*ast.UnaryExprOrTypeTraitExpr); ok {
			allocType = v.Type2
		} else {
			return nil, "", nil, nil,
				fmt.Errorf("Unsupport type '%T' in function calloc", n.Children()[2])
		}
		goType, err := types.ResolveType(p, allocType)
		if err != nil {
			return nil, "", nil, nil, err
		}
		return &goast.CallExpr{
			Fun: util.NewIdent("make"),
			Args: []goast.Expr{
				&goast.ArrayType{Elt: goast.NewIdent(goType)},
				size,
			},
		}, allocType + " *", preStmts, postStmts, nil
	}

	// Get the function definition from it's name. The case where it is not
	// defined is handled below (we haven't seen the prototype yet).
	functionDef := program.GetFunctionDefinition(functionName)

	if functionDef == nil {
		// We do not have a prototype for the function, but we should not exit
		// here. Instead we will create a mock definition for it so that this
		// transpile function will always return something and continue.
		//
		// The mock function definition is never actually saved to the program
		// definitions, so each time we see the CallExpr it will run this every
		// time. This is so if we come across the real prototype later it will
		// be handled correctly. Or at least "more" correctly.
		functionDef = &program.FunctionDefinition{
			Name: functionName,
		}
		if len(n.Children()) > 0 {
			if v, ok := n.Children()[0].(*ast.ImplicitCastExpr); ok && (types.IsFunction(v.Type) || types.IsTypedefFunction(p, v.Type)) {
				t := v.Type
				if types.IsTypedefFunction(p, t) {
					t = t[0 : len(t)-len(" *")]
					t, _ = p.TypedefType[t]
				}
				fields, returns, err := types.ParseFunction(t)
				if err != nil {
					p.AddMessage(p.GenerateWarningMessage(fmt.Errorf("Cannot resolve function : %v", err), n))
					return nil, "", nil, nil, err
				}
				functionDef.ReturnType = returns[0]
				functionDef.ArgumentTypes = fields
			}
		}
	} else {
		// type correction for definition function in
		// package program
		var ok bool
		for pos, arg := range n.Children() {
			if pos == 0 {
				continue
			}
			if pos >= len(functionDef.ArgumentTypes) {
				continue
			}
			if arg, ok = arg.(*ast.ImplicitCastExpr); ok {
				arg.(*ast.ImplicitCastExpr).Type = functionDef.ArgumentTypes[pos-1]
			}
		}
	}

	if functionDef.Substitution != "" {
		parts := strings.Split(functionDef.Substitution, ".")
		importName := strings.Join(parts[:len(parts)-1], ".")
		p.AddImport(importName)

		parts2 := strings.Split(functionDef.Substitution, "/")
		functionName = parts2[len(parts2)-1]
	}

	args := []goast.Expr{}
	argTypes := []string{}
	i := 0
	for _, arg := range n.Children()[1:] {
		e, eType, newPre, newPost, err := transpileToExpr(arg, p, false)
		if err != nil {
			return nil, "unknown2", nil, nil, err
		}
		argTypes = append(argTypes, eType)

		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

		_, arraySize := types.GetArrayTypeAndSize(eType)

		// If we are using varargs with Printf we need to make sure that certain
		// types are cast correctly.
		if functionName == "fmt.Printf" {
			// Make sure that any string parameters (const char*) are truncated
			// to the NULL byte.
			if arraySize != -1 {
				p.AddImport("github.com/elliotchance/c2go/noarch")
				e = util.NewCallExpr(
					"noarch.CStringToString",
					&goast.SliceExpr{X: e},
				)
			}

			// Byte slices (char*) must also be truncated to the NULL byte.
			//
			// TODO: This would also apply to other formatting functions like
			// fprintf, etc.
			if i > len(functionDef.ArgumentTypes)-1 &&
				(eType == "char *" || eType == "char*") {
				p.AddImport("github.com/elliotchance/c2go/noarch")
				e = util.NewCallExpr("noarch.CStringToString", e)
			}
		}

		args = append(args, e)

		i++
	}

	// These are the arguments once any transformations have taken place.
	realArgs := []goast.Expr{}

	// Apply transformation if needed. A transformation rearranges the return
	// value(s) and parameters. It is also used to indicate when a variable must
	// be passed by reference.
	if functionDef.ReturnParameters != nil || functionDef.Parameters != nil {
		for i, a := range functionDef.Parameters {
			byReference := false

			// Negative position means that it must be passed by reference.
			if a < 0 {
				byReference = true
				a = -a
			}

			// Rearrange the arguments. The -1 is because 0 would be the return
			// value.
			realArg := args[a-1]

			if byReference {
				// We have to create a temporary variable to pass by reference.
				// Then we can assign the real variable from it.
				realArg = &goast.UnaryExpr{
					Op: token.AND,
					X:  args[i],
				}
			} else {
				realArg, err = types.CastExpr(p, realArg, argTypes[i],
					functionDef.ArgumentTypes[i])
				p.AddMessage(
					p.GenerateWarningOrErrorMessage(err, n, realArg == nil),
				)

				if realArg == nil {
					realArg = util.NewNil()
				}
			}

			realArgs = append(realArgs, realArg)
		}
	} else {
		// Keep all the arguments the same. But make sure we cast to the correct
		// types.
		for i, a := range args {
			if i > len(functionDef.ArgumentTypes)-1 {
				// This means the argument is one of the varargs so we don't
				// know what type it needs to be cast to.
			} else {
				a, err = types.CastExpr(p, a, argTypes[i],
					functionDef.ArgumentTypes[i])

				if p.AddMessage(p.GenerateWarningMessage(err, n)) {
					a = util.NewNil()
				}
			}

			realArgs = append(realArgs, a)
		}
	}

	// Added for support removing function `free` of <stdlib.h>
	// Example of C code:
	// free(i+=4,buffer)
	// Example of result Go code:
	// i += 4
	// _ = buffer
	if functionDef.Substitution == "_" {
		devNull := &goast.AssignStmt{
			Lhs: []goast.Expr{goast.NewIdent("_")},
			Tok: token.ASSIGN,
			Rhs: []goast.Expr{realArgs[0]},
		}
		preStmts = append(preStmts, devNull)
		return nil, "", preStmts, postStmts, nil
	}

	return util.NewCallExpr(functionName, realArgs...),
		functionDef.ReturnType, preStmts, postStmts, nil
}
