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
