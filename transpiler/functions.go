// This file contains functions for declaring function prototypes, expressions
// that call functions, returning from function and the coordination of
// processing the function bodies.

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

// getFunctionBody returns the function body as a CompoundStmt. If the function
// is a prototype or forward declaration (meaning it has no body) then nil is
// returned.
func getFunctionBody(n *ast.FunctionDecl) *ast.CompoundStmt {
	// It's possible that the last node is the CompoundStmt (after all the
	// parameter declarations) - but I don't know this for certain so we will
	// look at all the children for now.
	for _, c := range n.Children() {
		if b, ok := c.(*ast.CompoundStmt); ok {
			return b
		}
	}

	return nil
}

// transpileFunctionDecl transpiles the function prototype.
//
// The function prototype may also have a body. If it does have a body the whole
// function will be transpiled into Go.
//
// If there is no function body we register the function interally (actually
// either way the function is registered internally) but we do not do anything
// because Go does not use or have any use for forward declarations of
// functions.
func transpileFunctionDecl(n *ast.FunctionDecl, p *program.Program) (
	decls []goast.Decl, err error) {
	var body *goast.BlockStmt

	// This is set at the start of the function declaration so when the
	// ReturnStmt comes alone it will know what the current function is, and
	// therefore be able to lookup what the real return type should be. I'm sure
	// there is a much better way of doing this.
	p.Function = n
	defer func() {
		// Reset the function name when we go out of scope.
		p.Function = nil
	}()

	n.Name = util.ConvertFunctionNameFromCtoGo(n.Name)

	// Always register the new function. Only from this point onwards will
	// we be allowed to refer to the function.
	if p.GetFunctionDefinition(n.Name) == nil {
		p.AddFunctionDefinition(program.FunctionDefinition{
			Name:          n.Name,
			ReturnType:    getFunctionReturnType(n.Type),
			ArgumentTypes: getFunctionArgumentTypes(n),
			Substitution:  "",
		})
	}

	// If the function has a direct substitute in Go we do not want to
	// output the C definition of it.
	f := p.GetFunctionDefinition(n.Name)
	if f != nil && f.Substitution != "" {
		err = nil
		return
	}

	// Test if the function has a body. This is identified by a child node that
	// is a CompoundStmt (since it is not valid to have a function body without
	// curly brackets).
	functionBody := getFunctionBody(n)
	if functionBody != nil {
		var pre, post []goast.Stmt
		body, pre, post, err = transpileToBlockStmt(functionBody, p)
		if err != nil || len(pre) > 0 || len(post) > 0 {
			p.AddMessage(p.GenerateErrorMessage(fmt.Errorf("Not correct result in function %s body: err = %v", n.Name, err), n))
			err = nil // Error is ignored
		}
	}

	// These functions cause us trouble for whatever reason. Some of them might
	// even work now.
	//
	// TODO: Some functions are ignored because they are too much trouble
	// https://github.com/elliotchance/c2go/issues/78
	if n.Name == "__istype" ||
		n.Name == "__isctype" ||
		n.Name == "__wcwidth" ||
		n.Name == "__sputc" ||
		n.Name == "__inline_signbitf" ||
		n.Name == "__inline_signbitd" ||
		n.Name == "__inline_signbitl" {
		err = nil
		return
	}

	if functionBody != nil {
		// If verbose mode is on we print the name of the function as a comment
		// immediately to stdout. This will appear at the top of the program but
		// make it much easier to diagnose when the transpiler errors.
		if p.Verbose {
			fmt.Printf("// Function: %s(%s)\n", f.Name,
				strings.Join(f.ArgumentTypes, ", "))
		}

		var fieldList = &goast.FieldList{}
		fieldList, err = getFieldList(n, p)
		if err != nil {
			return
		}

		t, err := types.ResolveType(p, f.ReturnType)
		p.AddMessage(p.GenerateWarningMessage(err, n))

		if p.Function != nil && p.Function.Name == "main" {
			// main() function does not have a return type.
			t = ""

			// This collects statements that will be placed at the top of
			// (before any other code) in main().
			prependStmtsInMain := []goast.Stmt{}

			// In Go, the main() function does not take the system arguments.
			// Instead they are accessed through the os package. We create new
			// variables in the main() function (if needed), immediately after
			// the __init() for these variables.
			if len(fieldList.List) > 0 {
				p.AddImport("os")

				prependStmtsInMain = append(
					prependStmtsInMain,
					&goast.AssignStmt{
						Lhs: []goast.Expr{fieldList.List[0].Names[0]},
						Tok: token.DEFINE,
						Rhs: []goast.Expr{util.NewCallExpr("int32", util.NewCallExpr("len", util.NewTypeIdent("os.Args")))},
					},
				)
			}

			if len(fieldList.List) > 1 {
				argvMultiArrayName := &goast.Ident{}
				argvArrayName := &goast.Ident{}
				*argvArrayName = *fieldList.List[1].Names[0]
				*argvMultiArrayName = *argvArrayName
				argvArrayName.Name += "__array"
				argvMultiArrayName.Name += "__multiarray"
				prependStmtsInMain = append(
					prependStmtsInMain,
					&goast.AssignStmt{
						Lhs: []goast.Expr{argvMultiArrayName},
						Tok: token.DEFINE,
						Rhs: []goast.Expr{&goast.CompositeLit{Type: util.NewTypeIdent("[][]byte")}},
					},
					&goast.AssignStmt{
						Lhs: []goast.Expr{argvArrayName},
						Tok: token.DEFINE,
						Rhs: []goast.Expr{&goast.CompositeLit{Type: util.NewTypeIdent("[]*byte")}},
					},
					&goast.RangeStmt{
						Key:   goast.NewIdent("_"),
						Value: util.NewIdent("argvSingle"),
						Tok:   token.DEFINE,
						X:     util.NewTypeIdent("os.Args"),
						Body: &goast.BlockStmt{
							List: []goast.Stmt{
								&goast.AssignStmt{
									Lhs: []goast.Expr{argvMultiArrayName},
									Tok: token.ASSIGN,
									Rhs: []goast.Expr{util.NewCallExpr(
										"append",
										argvMultiArrayName,
										util.NewCallExpr("append",
											util.NewCallExpr("[]byte", util.NewIdent("argvSingle")),
											util.NewIntLit(0)),
									)},
								},
							},
						},
					},
					&goast.RangeStmt{
						Key:   goast.NewIdent("_"),
						Value: util.NewIdent("argvSingle"),
						Tok:   token.DEFINE,
						X:     argvMultiArrayName,
						Body: &goast.BlockStmt{
							List: []goast.Stmt{
								&goast.AssignStmt{
									Lhs: []goast.Expr{argvArrayName},
									Tok: token.ASSIGN,
									Rhs: []goast.Expr{util.NewCallExpr(
										"append",
										argvArrayName,
										&goast.UnaryExpr{
											Op: token.AND,
											X: &goast.IndexExpr{
												X:     util.NewIdent("argvSingle"),
												Index: util.NewIntLit(0),
											},
										},
									)},
								},
							},
						},
					},
					&goast.AssignStmt{
						Lhs: []goast.Expr{fieldList.List[1].Names[0]},
						Tok: token.DEFINE,
						Rhs: []goast.Expr{
							&goast.StarExpr{
								X: &goast.CallExpr{
									Fun: &goast.ParenExpr{
										X: util.NewTypeIdent("***byte"),
									},
									Args: []goast.Expr{
										util.NewCallExpr("unsafe.Pointer", &goast.UnaryExpr{
											Op: token.AND,
											X:  argvArrayName,
										}),
									},
								},
							},
						},
					})
			}

			// Prepend statements for main().
			body.List = append(prependStmtsInMain, body.List...)

			// The main() function does not have arguments or a return value.
			fieldList = &goast.FieldList{}
		}

		// Each function MUST have "ReturnStmt",
		// except function without return type
		var addReturnName bool
		if len(body.List) > 0 {
			last := body.List[len(body.List)-1]
			if _, ok := last.(*goast.ReturnStmt); !ok && t != "" {
				body.List = append(body.List, &goast.ReturnStmt{})
				addReturnName = true
			}
		}

		decls = append(decls, &goast.FuncDecl{
			Name: util.NewIdent(n.Name),
			Type: util.NewFuncType(fieldList, t, addReturnName),
			Body: body,
		})
	}

	err = nil
	return
}

// getFieldList returns the parameters of a C function as a Go AST FieldList.
func getFieldList(f *ast.FunctionDecl, p *program.Program) (_ *goast.FieldList, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Error in function field list. err = %v", err)
		}
	}()
	r := []*goast.Field{}
	for _, n := range f.Children() {
		if v, ok := n.(*ast.ParmVarDecl); ok {
			if types.IsFunction(v.Type) {
				field, err := newFunctionField(p, v.Name, v.Type)
				if err != nil {
					p.AddMessage(p.GenerateWarningMessage(err, v))
					continue
				}
				r = append(r, field)
				continue
			}
			// when passing va_list to a function, always name it c2goVaList
			if v.Type == "struct __va_list_tag *" {
				v.Name = "c2goVaList"
			}
			t, err := types.ResolveType(p, v.Type)
			p.AddMessage(p.GenerateWarningMessage(err, f))

			r = append(r, &goast.Field{
				Names: []*goast.Ident{util.NewIdent(v.Name)},
				Type:  util.NewTypeIdent(t),
			})
		}
	}

	// for function argument: ...
	if strings.Contains(f.Type, "...") {
		r = append(r, &goast.Field{
			Names: []*goast.Ident{util.NewIdent("c2goArgs")},
			Type: &goast.Ellipsis{
				Ellipsis: 1,
				Elt: &goast.InterfaceType{
					Interface: 1,
					Methods: &goast.FieldList{
						Opening: 1,
					},
					Incomplete: false,
				},
			},
		})
	}

	return &goast.FieldList{
		List: r,
	}, nil
}

func transpileReturnStmt(n *ast.ReturnStmt, p *program.Program) (
	_ goast.Stmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileReturnStmt. err = %v", err)
		}
	}()
	// There may not be a return value. Then we don't have to both ourselves
	// with all the rest of the logic below.
	if len(n.Children()) == 0 {
		return &goast.ReturnStmt{}, nil, nil, nil
	}

	var eType string
	var e goast.Expr
	e, eType, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, nil, nil, err
	}
	if e == nil {
		return nil, nil, nil, fmt.Errorf("Expr is nil")
	}

	f := p.GetFunctionDefinition(p.Function.Name)

	t, err := types.CastExpr(p, e, eType, f.ReturnType)
	if p.AddMessage(p.GenerateWarningMessage(err, n)) {
		t = util.NewNil()
	}

	results := []goast.Expr{t}

	// main() function is not allowed to return a result. Use os.Exit if
	// non-zero.
	if p.Function != nil && p.Function.Name == "main" {
		litExpr, isLiteral := getReturnLiteral(e)
		if !isLiteral || (isLiteral && litExpr.Value != "0") {
			p.AddImport("os")
			return util.NewExprStmt(util.NewCallExpr("os.Exit", util.NewCallExpr("int", results...))),
				preStmts, postStmts, nil
		}
		results = []goast.Expr{}
	}

	return &goast.ReturnStmt{
		Results: results,
	}, preStmts, postStmts, nil
}

func getReturnLiteral(e goast.Expr) (litExpr *goast.BasicLit, ok bool) {
	if litExpr, ok = e.(*goast.BasicLit); ok {
		return
	}
	if callExpr, ok2 := e.(*goast.CallExpr); ok2 {
		if funExpr, ok3 := callExpr.Fun.(*goast.Ident); !ok3 || funExpr.Name != "int32" {
			return nil, false
		}
		if len(callExpr.Args) != 1 {
			return nil, false
		}
		if litExpr, ok = callExpr.Args[0].(*goast.BasicLit); ok {
			return
		}
	}
	return nil, false
}

func getFunctionReturnType(f string) string {
	// The C type of the function will be the complete prototype, like:
	//
	//     __inline_isfinitef(float) int
	//
	// will have a C type of:
	//
	//     int (float)
	//
	// The arguments will handle themselves, we only care about the return type
	// ('int' in this case)
	returnType := strings.TrimSpace(strings.Split(f, "(")[0])

	if returnType == "" {
		panic(fmt.Sprintf("unable to extract the return type from: %s", f))
	}

	return returnType
}

// getFunctionArgumentTypes returns the C types of the arguments in a function.
func getFunctionArgumentTypes(f *ast.FunctionDecl) []string {
	r := []string{}
	for _, n := range f.Children() {
		if v, ok := n.(*ast.ParmVarDecl); ok {
			r = append(r, v.Type)
		}
	}

	return r
}
