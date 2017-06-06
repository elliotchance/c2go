// This file contains functions for transpiling unary operator expressions.

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

func transpileUnaryOperator(n *ast.UnaryOperator, p *program.Program) (
    goast.Expr, string, []goast.Stmt, []goast.Stmt, error) {
    preStmts := []goast.Stmt{}
    postStmts := []goast.Stmt{}
    operator := getTokenForOperator(n.Operator)

    // Unfortunately we cannot use the Go increment operators because we are not
    // providing any position information for tokens. This means that the ++/--
    // would be placed before the expression and would be invalid in Go.
    //
    // Until it can be properly fixed (can we trick Go into to placing it after
    // the expression with a magic position?) we will have to return a
    // BinaryExpr with the same functionality.
    if operator == token.INC || operator == token.DEC {
        // Construct code for assigning value to an union field
        memberExpr, ok := n.Children[0].(*ast.MemberExpr)
        if ok {
            ref := memberExpr.GetDeclRef()
            if ref != nil {
                typename, err := types.ResolveType(p, ref.Type)
                if err != nil {
                    return nil, "", preStmts, postStmts, err
                }

                if typename[0] == '*' {
                    typename = typename[1:]
                }

                binaryOperator := token.ADD
                if operator == token.DEC {
                    binaryOperator = token.SUB
                }

                union := p.GetStruct(typename)
                if union.IsUnion {
                    // Method suffix for using getters and setters of Go union type
                    methodSuffix := strings.Title(memberExpr.Name)

                    // Method names
                    getterName := fmt.Sprintf("%s.Get%s", ref.Name, methodSuffix)
                    setterName := fmt.Sprintf("%s.Set%s", ref.Name, methodSuffix)

                    // Call-Expression argument
                    argLhs := util.NewCallExpr(getterName)
                    argOp := getTokenForOperator(binaryOperation)
                    argRhs := util.NewIntLit(1)
                    argValue := util.NewBinaryExpr(argLhs, argOp, argRhs)

                    // Make Go expression
                    resExpr := util.NewCallExpr(funcName, argValue)

                    return resExpr, n.Type, preStmts, postStmts, nil
                }
            }
        }

        binaryOperator := "+="
        if operator == token.DEC {
            binaryOperator = "-="
        }

        return transpileBinaryOperator(&ast.BinaryOperator{
            Type:     n.Type,
            Operator: binaryOperator,
            Children: []ast.Node{
                n.Children[0], &ast.IntegerLiteral{
                    Type:     "int",
                    Value:    "1",
                    Children: []ast.Node{},
                },
            },
        }, p)
    }

    // Otherwise handle like a unary operator.
    e, eType, newPre, newPost, err := transpileToExpr(n.Children[0], p)
    if err != nil {
        return nil, "", nil, nil, err
    }

    preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

    if operator == token.NOT {
        if eType == "bool" || eType == "_Bool" {
            return &goast.UnaryExpr{
                X:  e,
                Op: operator,
            }, "bool", preStmts, postStmts, nil
        }

        t, err := types.ResolveType(p, eType)
        p.AddMessage(ast.GenerateWarningMessage(err, n))

        if t == "[]byte" {
            return util.NewUnaryExpr(
                token.NOT, util.NewCallExpr("noarch.CStringIsNull", e),
            ), "bool", preStmts, postStmts, nil
        }

        p.AddImport("github.com/elliotchance/c2go/noarch")

        functionName := fmt.Sprintf("noarch.Not%s", util.Ucfirst(t))

        return util.NewCallExpr(functionName, e),
            eType, preStmts, postStmts, nil
    }

    // Dereferencing.
    if operator == token.MUL {
        if eType == "const char *" {
            return &goast.IndexExpr{
                X:     e,
                Index: util.NewIntLit(0),
            }, "char", preStmts, postStmts, nil
        }

        t, err := types.GetDereferenceType(eType)
        if err != nil {
            return nil, "", preStmts, postStmts, err
        }

        // C is more relaxed with this syntax. In Go we convert all of the
        // pointers to slices, so we have to be careful when dereference a slice
        // that it actually takes the first element instead.
        resolvedType, err := types.ResolveType(p, eType)
        if strings.HasPrefix(resolvedType, "[]") {
            return &goast.IndexExpr{
                X:     e,
                Index: util.NewIntLit(0),
            }, t, preStmts, postStmts, nil
        }

        return &goast.StarExpr{
            X: e,
        }, t, preStmts, postStmts, nil
    }

    if operator == token.AND {
        // We now have a pointer to the original type.
        eType += " *"
    }

    return &goast.UnaryExpr{
        Op: operator,
        X:  e,
    }, eType, preStmts, postStmts, nil
}

func transpileUnaryExprOrTypeTraitExpr(n *ast.UnaryExprOrTypeTraitExpr, p *program.Program) (
    *goast.BasicLit, string, []goast.Stmt, []goast.Stmt, error) {
    t := n.Type2

    // It will have children if the sizeof() is referencing a variable.
    // Fortunately clang already has the type in the AST for us.
    if len(n.Children) > 0 {
        var realFirstChild interface{}
        t = ""

        switch c := n.Children[0].(type) {
        case *ast.ParenExpr:
            realFirstChild = c.Children[0]
        case *ast.DeclRefExpr:
            t = c.Type
        default:
            panic(fmt.Sprintf("cannot find first child from: %#v", n.Children[0]))
        }

        if t == "" {
            switch ty := realFirstChild.(type) {
            case *ast.DeclRefExpr:
                t = ty.Type2

            case *ast.ArraySubscriptExpr:
                t = ty.Type

            case *ast.MemberExpr:
                t = ty.Type

            case *ast.UnaryOperator:
                t = ty.Type

            case *ast.ParenExpr:
                t = ty.Type

            case *ast.CallExpr:
                t = ty.Type

            default:
                panic(fmt.Sprintf("cannot do unary on: %#v", ty))
            }
        }
    }

    sizeInBytes, err := types.SizeOf(p, t)
    p.AddMessage(ast.GenerateWarningMessage(err, n))

    return util.NewIntLit(sizeInBytes), n.Type1, nil, nil, nil
}
