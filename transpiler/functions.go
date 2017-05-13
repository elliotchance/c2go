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

	goast "go/ast"
)

// transpileFunctionDecl transpiles the function prototype.
//
// The function prototype may also have a body. If it does have a body the whole
// function will be transpiled into Go.
//
// If there is no function body we register the function interally (actually
// either way the function is registered internally) but we do not do anything
// becuase Go does not use or have any use for forward declarations of
// functions.
func transpileFunctionDecl(n *ast.FunctionDecl, p *program.Program) error {
	var body *goast.BlockStmt

	// This is set at the start of the function declaration so when the
	// ReturnStmt comes alone it will know what the current function is, and
	// therefore be able to lookup what the real return type should be. I'm sure
	// there is a much better way of doing this.
	p.FunctionName = n.Name
	defer func() {
		// Reset the function name when we go out of scope.
		p.FunctionName = ""
	}()

	// Always register the new function. Only from this point onwards will
	// we be allowed to refer to the function.
	if program.GetFunctionDefinition(n.Name) == nil {
		program.AddFunctionDefinition(program.FunctionDefinition{
			Name:          n.Name,
			ReturnType:    getFunctionReturnType(n.Type),
			ArgumentTypes: getFunctionArgumentTypes(n),
			Substitution:  "",
		})
	}

	// If the function has a direct substitute in Go we do not want to
	// output the C definition of it.
	f := program.GetFunctionDefinition(n.Name)
	if f != nil && f.Substitution != "" {
		return nil
	}

	// Test if the function has a body. This is identified by a child node that
	// is a CompoundStmt (since it is not valid to have a function body without
	// curly brackets).
	//
	// It's possible that the last node is the CompoundStmt (after all the
	// parameter declarations) - but I don't know this for certain so we will
	// look at all the children for now.
	hasBody := false
	for _, c := range n.Children {
		if b, ok := c.(*ast.CompoundStmt); ok {
			var err error
			var newPre, newPost []goast.Stmt

			body, newPre, newPost, err = transpileToBlockStmt(b, p)
			if err != nil {
				return err
			}

			if len(newPre) > 0 || len(newPost) > 0 {
				panic("bad")
			}

			hasBody = true
			break
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
		return nil
	}

	if hasBody {
		fieldList, err := getFieldList(n, p)
		if err != nil {
			return err
		}

		returnTypes := []*goast.Field{
			&goast.Field{
				Type: goast.NewIdent(types.ResolveType(p, f.ReturnType)),
			},
		}

		if p.FunctionName == "main" {
			// main() function does not have a return type.
			returnTypes = []*goast.Field{}

			// We also need to append a setup function that will instantiate
			// some things that are expected to be available at runtime.
			body.List = append([]goast.Stmt{
				&goast.ExprStmt{
					X: &goast.CallExpr{
						Fun: goast.NewIdent("__init"),
					},
				},
			}, body.List...)
		}

		p.File.Decls = append(p.File.Decls, &goast.FuncDecl{
			Name: goast.NewIdent(n.Name),
			Type: &goast.FuncType{
				Params: fieldList,
				Results: &goast.FieldList{
					List: returnTypes,
				},
			},
			Body: body,
		})
	}

	return nil
}

// getFieldList returns the paramaters of a C function as a Go AST FieldList.
func getFieldList(f *ast.FunctionDecl, p *program.Program) (*goast.FieldList, error) {
	// The main() function does not have arguments or a return value.
	if f.Name == "main" {
		return &goast.FieldList{}, nil
	}

	r := []*goast.Field{}
	for _, n := range f.Children {
		if v, ok := n.(*ast.ParmVarDecl); ok {
			r = append(r, &goast.Field{
				Names: []*goast.Ident{goast.NewIdent(v.Name)},
				Type:  goast.NewIdent(types.ResolveType(p, v.Type)),
			})
		}
	}

	return &goast.FieldList{
		List: r,
	}, nil
}

func transpileReturnStmt(n *ast.ReturnStmt, p *program.Program) (
	goast.Stmt, []goast.Stmt, []goast.Stmt, error) {
	// There may not be a return value. Then we don't have to both ourselves
	// with all the rest of the logic below.
	if len(n.Children) == 0 {
		return &goast.ReturnStmt{}, nil, nil, nil
	}

	e, eType, preStmts, postStmts, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, nil, nil, err
	}

	f := program.GetFunctionDefinition(p.FunctionName)

	results := []goast.Expr{types.CastExpr(p, e, eType, f.ReturnType)}

	// main() function is not allowed to return a result. Use os.Exit if non-zero
	if p.FunctionName == "main" {
		litExpr, isLiteral := e.(*goast.BasicLit)
		if !isLiteral || (isLiteral && litExpr.Value != "0") {
			p.AddImport("os")
			return &goast.ExprStmt{
				X: &goast.CallExpr{
					Fun:  goast.NewIdent("os.Exit"),
					Args: results,
				},
			}, preStmts, postStmts, nil
		}
		results = []goast.Expr{}
	}

	return &goast.ReturnStmt{
		Results: results,
	}, preStmts, postStmts, nil
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
	for _, n := range f.Children {
		if v, ok := n.(*ast.ParmVarDecl); ok {
			r = append(r, v.Type)
		}
	}

	return r
}
