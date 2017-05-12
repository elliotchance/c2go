// This file contains functions for declaring function prototypes, expressions
// that call functions, returning from function and the coordination of
// processing the function bodies.

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

	// The first child will always contain the name of the function being
	// called.
	firstChild := n.Children[0].(*ast.ImplicitCastExpr).Children[0]
	functionName := firstChild.(*ast.DeclRefExpr).Name

	// Get the function definition from it's name. The case where it is not
	// defined is handled below (we haven't seen the prototype yet).
	functionDef := program.GetFunctionDefinition(functionName)

	if functionDef == nil {
		errorMessage := fmt.Sprintf("unknown function: %s", functionName)
		return nil, "", nil, nil, errors.New(errorMessage)
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

		// if i > len(functionDef.ArgumentTypes)-1 {
		// 	// This means the argument is one of the varargs so we don't know
		// 	// what type it needs to be cast to.
		// } else {
		// 	//e = types.CastExpr(p, e, eType, functionDef.ArgumentTypes[i])
		// }

		_, arraySize := types.GetArrayTypeAndSize(eType)
		if functionName == "fmt.Printf" && arraySize != -1 {
			p.AddImport("github.com/elliotchance/c2go/noarch")
			e = util.NewCallExpr(
				"noarch.NullTerminatedString",
				util.NewCallExpr("string", &goast.SliceExpr{X: e}),
			)
		}

		// We cannot use preallocated byte slices as strings in the same way we
		// can do it in C. Instead we have to create a temporary string
		// variable.
		if functionName == "noarch.Fscanf" && arraySize != -1 {
			// FIXME: The name of the temp variable needs to be random.

			// var __temp string
			preStmts = append(preStmts, &goast.DeclStmt{
				&goast.GenDecl{
					Tok: token.VAR,
					Specs: []goast.Spec{
						&goast.ValueSpec{
							Names: []*goast.Ident{goast.NewIdent("__temp")},
							Type:  goast.NewIdent("string"),
						},
					},
				},
			})

			postStmts = append(postStmts, &goast.ExprStmt{
				X: util.NewCallExpr("copy", &goast.SliceExpr{
					X: e,
				}, goast.NewIdent("__temp")),
			})

			e = &goast.UnaryExpr{
				Op: token.AND,
				X:  goast.NewIdent("__temp"),
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
				realArg = types.CastExpr(p, realArg, argTypes[i], functionDef.ArgumentTypes[i])
			}

			realArgs = append(realArgs, realArg)
		}
	} else {
		// Keep all the arguments the same. But make sure we cast to the correct
		// types.
		for i, a := range args {
			if i > len(functionDef.ArgumentTypes)-1 {
				// This means the argument is one of the varargs so we don't know
				// what type it needs to be cast to.
			} else {
				a = types.CastExpr(p, a, argTypes[i], functionDef.ArgumentTypes[i])
			}

			realArgs = append(realArgs, a)
		}
	}

	return &goast.CallExpr{
		Fun:  goast.NewIdent(functionName),
		Args: realArgs,
	}, functionDef.ReturnType, preStmts, postStmts, nil
}

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
