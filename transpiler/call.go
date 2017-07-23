// This file contains functions for transpiling function calls (invocations).

package transpiler

import (
	"errors"
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
		return fc.Name

	case *ast.ParenExpr:
		return getName(fc.Children[0])

	case *ast.UnaryOperator:
		return getName(fc.Children[0])

	case *ast.ImplicitCastExpr:
		return getName(fc.Children[0])

	case *ast.CStyleCastExpr:
		return getName(fc.Children[0])

	default:
		panic(fmt.Sprintf("cannot CallExpr on: %#v", fc))
	}
}

func getNameOfFunctionFromCallExpr(n *ast.CallExpr) (string, error) {
	// The first child will always contain the name of the function being
	// called.
	firstChild, ok := n.Children[0].(*ast.ImplicitCastExpr)
	if !ok {
		err := fmt.Errorf("unable to use CallExpr: %#v", n.Children[0])
		return "", err
	}

	return getName(firstChild.Children[0]), nil
}

// transpileCallExpr transpiles expressions that calls a function, for example:
//
//     foo("bar")
//
// It returns three arguments; the Go AST expression, the C type (that is
// returned by the function) and any error. If there is an error returned you
// can assume the first two arguments will not contain any useful information.
func transpileCallExpr(n *ast.CallExpr, p *program.Program) (
	*goast.CallExpr, string, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	functionName, err := getNameOfFunctionFromCallExpr(n)
	if err != nil {
		return nil, "", nil, nil, err
	}

	// Get the function definition from it's name. The case where it is not
	// defined is handled below (we haven't seen the prototype yet).
	functionDef := program.GetFunctionDefinition(functionName)

	if functionDef == nil {
		errorMessage := fmt.Sprintf("unknown function: %s", functionName)
		p.AddMessage(ast.GenerateWarningMessage(errors.New(errorMessage), n))

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
	for _, arg := range n.Children[1:] {
		e, eType, newPre, newPost, err := transpileToExpr(arg, p)
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
					ast.GenerateWarningOrErrorMessage(err, n, realArg == nil),
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

				if p.AddMessage(ast.GenerateWarningMessage(err, n)) {
					a = util.NewNil()
				}
			}

			realArgs = append(realArgs, a)
		}
	}

	return util.NewCallExpr(functionName, realArgs...),
		functionDef.ReturnType, preStmts, postStmts, nil
}
