// This file contains functions for transpiling function calls (invocations).

package transpiler

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"

	goast "go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

func getMemberName(firstChild ast.Node) (name string, ok bool) {
	switch fc := firstChild.(type) {
	case *ast.MemberExpr:
		return fc.Name, true

	case *ast.ParenExpr:
		return getMemberName(fc.Children()[0])

	case *ast.ImplicitCastExpr:
		return getMemberName(fc.Children()[0])

	case *ast.CStyleCastExpr:
		return getMemberName(fc.Children()[0])

	}
	return "", false
}

func getName(p *program.Program, firstChild ast.Node) (name string, err error) {
	switch fc := firstChild.(type) {
	case *ast.DeclRefExpr:
		return fc.Name, nil

	case *ast.MemberExpr:
		if isUnionMemberExpr(p, fc) {
			var expr goast.Expr
			expr, _, _, _, err = transpileToExpr(fc, p, false)
			if err != nil {
				return
			}
			var buf bytes.Buffer
			err = printer.Fprint(&buf, token.NewFileSet(), expr)
			if err != nil {
				return
			}
			return buf.String(), nil
		}
		if len(fc.Children()) == 0 {
			return fc.Name, nil
		}
		var n string
		n, err = getName(p, fc.Children()[0])
		if err != nil {
			return
		}
		return n + "." + fc.Name, nil

	case *ast.ParenExpr:
		return getName(p, fc.Children()[0])

	case *ast.UnaryOperator:
		return getName(p, fc.Children()[0])

	case *ast.ImplicitCastExpr:
		return getName(p, fc.Children()[0])

	case *ast.CStyleCastExpr:
		return getName(p, fc.Children()[0])

	case *ast.ArraySubscriptExpr:
		var expr goast.Expr
		expr, _, _, _, err = transpileArraySubscriptExpr(fc, p, false)
		if err != nil {
			return
		}
		var buf bytes.Buffer
		err = printer.Fprint(&buf, token.NewFileSet(), expr)
		if err != nil {
			return
		}
		return buf.String(), nil
	}

	return "", fmt.Errorf("cannot getName for: %#v", firstChild)
}

func getNameOfFunctionFromCallExpr(p *program.Program, n *ast.CallExpr) (string, error) {
	// The first child will always contain the name of the function being
	// called.
	firstChild, ok := n.Children()[0].(*ast.ImplicitCastExpr)
	if !ok {
		err := fmt.Errorf("unable to use CallExpr: %#v", n.Children()[0])
		return "", err
	}

	return getName(p, firstChild.Children()[0])
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
	defer func() {
		if err != nil {
			err = fmt.Errorf("Error in transpileCallExpr : %v", err)
		}
	}()

	functionName, err := getNameOfFunctionFromCallExpr(p, n)
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

	// function "calloc" from stdlib.c
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

	// function "qsort" from stdlib.h
	if functionName == "qsort" && len(n.Children()) == 5 {
		defer func() {
			if err != nil {
				err = fmt.Errorf("Function: qsort. err = %v", err)
			}
		}()
		/*
			CallExpr 0x2c6b1b0 <line:182:2, col:40> 'void'
			|-ImplicitCastExpr 0x2c6b198 <col:2> 'void (*)(void *, size_t, size_t, __compar_fn_t)' <FunctionToPointerDecay>
			| `-DeclRefExpr 0x2c6b070 <col:2> 'void (void *, size_t, size_t, __compar_fn_t)' Function 0x2bec110 'qsort' 'void (void *, size_t, size_t, __compar_fn_t)'
			|-ImplicitCastExpr 0x2c6b210 <col:9> 'void *' <BitCast>
			| `-ImplicitCastExpr 0x2c6b1f8 <col:9> 'int *' <ArrayToPointerDecay>
			|   `-DeclRefExpr 0x2c6b098 <col:9> 'int [6]' lvalue Var 0x2c6a6c0 'values' 'int [6]'
			|-ImplicitCastExpr 0x2c6b228 <col:17> 'size_t':'unsigned long' <IntegralCast>
			| `-IntegerLiteral 0x2c6b0c0 <col:17> 'int' 6
			|-UnaryExprOrTypeTraitExpr 0x2c6b0f8 <col:20, col:30> 'unsigned long' sizeof 'int'
			`-ImplicitCastExpr 0x2c6b240 <col:33> 'int (*)(const void *, const void *)' <FunctionToPointerDecay>
			  `-DeclRefExpr 0x2c6b118 <col:33> 'int (const void *, const void *)' Function 0x2c6aa70 'compare' 'int (const void *, const void *)'
		*/
		var element [4]goast.Expr
		for i := 1; i < 5; i++ {
			el, _, newPre, newPost, err := transpileToExpr(n.Children()[i], p, false)
			if err != nil {
				return nil, "", nil, nil, err
			}
			element[i-1] = el
			preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		}
		// found the C type
		t := n.Children()[3].(*ast.UnaryExprOrTypeTraitExpr).Type2
		t, err := types.ResolveType(p, t)
		if err != nil {
			return nil, "", nil, nil, err
		}

		var compareFunc string
		if v, ok := element[3].(*goast.Ident); ok {
			compareFunc = v.Name
		} else {
			return nil, "", nil, nil,
				fmt.Errorf("golang ast for compare function have type %T, expect ast.Ident", element[3])
		}

		var varName string
		id := extractArray(element[0])
		if id != nil {
			varName = id.Name
		} else {
			return nil, "", nil, nil, fmt.Errorf("cannot determine variable to be sorted")
		}

		p.AddImport("sort")
		src := fmt.Sprintf(`package main
		var %s func(a,b interface{})int
		var temp = func(i, j int) bool {
			c2goTempVarA := unsafe.Pointer(&%s[i])
			c2goTempVarB := unsafe.Pointer(&%s[j])
			return %s(c2goTempVarA, c2goTempVarB) <= 0
		}`, compareFunc, varName, varName, compareFunc)

		// Create the AST by parsing src.
		fset := token.NewFileSet() // positions are relative to fset
		f, err := parser.ParseFile(fset, "", src, 0)
		if err != nil {
			return nil, "", nil, nil, err
		}

		// AST tree part of code after "var temp = ..."
		convertExpr := f.Decls[1].(*goast.GenDecl).Specs[0].(*goast.ValueSpec).Values[0]

		return &goast.CallExpr{
			Fun: &goast.SelectorExpr{
				X:   goast.NewIdent("sort"),
				Sel: goast.NewIdent("SliceStable"),
			},
			Args: []goast.Expr{
				id,
				convertExpr,
			},
		}, "", preStmts, postStmts, nil
	}

	// Get the function definition from it's name. The case where it is not
	// defined is handled below (we haven't seen the prototype yet).
	functionDef := p.GetFunctionDefinition(functionName)

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
				if v, ok := p.TypedefType[t]; ok {
					t = v
				} else {
					if types.IsTypedefFunction(p, t) {
						t = t[0 : len(t)-len(" *")]
						t, _ = p.TypedefType[t]
					}
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

			if realArg == nil {
				return nil, "", preStmts, postStmts,
					fmt.Errorf("Real argument is nil in function : %s", functionName)
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

			if a == nil {
				return nil, "", preStmts, postStmts,
					fmt.Errorf("Argument is nil in function : %s", functionName)
			}

			if len(functionDef.ArgumentTypes) > i {
				if !types.IsPointer(p, functionDef.ArgumentTypes[i]) {
					if strings.HasPrefix(functionDef.ArgumentTypes[i], "union ") {
						a = &goast.CallExpr{
							Fun: &goast.SelectorExpr{
								X:   a,
								Sel: goast.NewIdent("copy"),
							},
							Lparen: 1,
						}
					}
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

func extractArray(expr goast.Expr) *goast.Ident {
	if v, ok := expr.(*goast.Ident); ok {
		return v
	}
	if se, ok := expr.(*goast.SliceExpr); ok {
		if v, ok2 := se.X.(*goast.Ident); ok2 {
			return v
		}
	}
	if ce, ok2 := expr.(*goast.CallExpr); ok2 && len(ce.Args) == 1 {
		if fid, ok3 := ce.Fun.(*goast.Ident); !ok3 || fid.Name != "unsafe.Pointer" {
			return nil
		}
		if ue, ok3 := ce.Args[0].(*goast.UnaryExpr); !ok3 || ue.Op != token.AND {
			return nil
		} else if idx, ok4 := ue.X.(*goast.IndexExpr); ok4 {
			if index, ok5 := idx.Index.(*goast.BasicLit); !ok5 || index.Value != "0" {
				return nil
			}
			if id, ok5 := idx.X.(*goast.Ident); ok5 {
				return id
			}
		}
	}
	return nil
}
