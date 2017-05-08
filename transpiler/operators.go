// This file contains functions transpiling unary and binary operator
// expressions.

package transpiler

import (
	"fmt"
	"go/token"

	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
)

func transpileBinaryOperator(n *ast.BinaryOperator, p *program.Program) (*goast.BinaryExpr, string, error) {
	left, leftType, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", err
	}

	right, rightType, err := transpileToExpr(n.Children[1], p)
	if err != nil {
		return nil, "", err
	}

	operator := getTokenForOperator(n.Operator)
	returnType := types.ResolveTypeForBinaryOperator(p, n.Operator, leftType, rightType)

	if operator == token.LAND {
		left = types.CastExpr(p, left, leftType, "bool")
		right = types.CastExpr(p, right, rightType, "bool")

		return &goast.BinaryExpr{
			X:  left,
			Op: operator,
			Y:  right,
		}, "bool", nil
	}

	// Convert "(0)" to "nil" when we are dealing with equality.
	if (operator == token.NEQ || operator == token.EQL) &&
		types.IsNullExpr(right) {
		if types.ResolveType(p, leftType) == "string" {
			right = &goast.BasicLit{
				Kind:  token.STRING,
				Value: `""`,
			}
		} else {
			right = goast.NewIdent("nil")
		}
	}

	if operator == token.ASSIGN {
		right = types.CastExpr(p, right, rightType, returnType)
	}

	return &goast.BinaryExpr{
		X:  left,
		Op: operator,
		Y:  right,
	}, types.ResolveTypeForBinaryOperator(p, n.Operator, leftType, rightType), nil
}

func transpileUnaryOperator(n *ast.UnaryOperator, p *program.Program) (goast.Expr, string, error) {
	operator := getTokenForOperator(n.Operator)

	// Unfortunately we cannot use the Go increment operators because we are not
	// providing any position information for tokens. This means that the ++/--
	// would be placed before the expression and would be invalid in Go.
	//
	// Until it can be properly fixed (can we trick Go into to placing it after
	// the expression with a magic position?) we will have to return a
	// BinaryExpr with the same functionality.
	if operator == token.INC || operator == token.DEC {
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
					Value:    1,
					Children: []ast.Node{},
				},
			},
		}, p)
	}

	// Otherwise handle like a unary operator.
	e, eType, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", err
	}

	if operator == token.NOT {
		if eType == "bool" || eType == "_Bool" {
			return &goast.UnaryExpr{
				X:  e,
				Op: operator,
			}, "bool", nil
		}

		t := types.ResolveType(p, eType)
		if t == "string" {
			return &goast.BinaryExpr{
				X:  e,
				Op: token.EQL,
				Y: &goast.BasicLit{
					Kind:  token.STRING,
					Value: `""`,
				},
			}, "bool", nil
		}

		p.AddImport("github.com/elliotchance/c2go/noarch")

		functionName := fmt.Sprintf("noarch.Not%s", util.Ucfirst(t))

		return &goast.CallExpr{
			Fun:  goast.NewIdent(functionName),
			Args: []goast.Expr{e},
		}, eType, nil
	}

	if operator == token.MUL {
		if eType == "const char *" {
			return &goast.IndexExpr{
				X: e,
				Index: &goast.BasicLit{
					Kind:  token.INT,
					Value: "0",
				},
			}, "char", nil
		}

		t, err := types.GetDereferenceType(eType)
		if err != nil {
			return nil, "", err
		}

		return &goast.StarExpr{
			X: e,
		}, t, nil
	}

	if operator == token.AND {
		eType += " *"
	}

	return &goast.UnaryExpr{
		X:  e,
		Op: operator,
	}, eType, nil
}

// transpileConditionalOperator transpiles a conditional (also known as a
// ternary) operator:
//
//     a ? b : c
//
// We cannot simply convert these to an "if" statement becuase they by inside
// another expression.
//
// Since Go does not support the ternary operator or inline "if" statements we
// use a function, noarch.Ternary() to work the same way.
//
// It is also important to note that C only evaulates the "b" or "c" condition
// based on the result of "a" (from the above example). So we wrap the "b" and
// "c" in closures so that the Ternary function will only evaluate one of them.
func transpileConditionalOperator(n *ast.ConditionalOperator, p *program.Program) (*goast.CallExpr, string, error) {
	a, _, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", err
	}

	b, _, err := transpileToExpr(n.Children[1], p)
	if err != nil {
		return nil, "", err
	}

	c, _, err := transpileToExpr(n.Children[2], p)
	if err != nil {
		return nil, "", err
	}

	p.AddImport("github.com/elliotchance/c2go/noarch")

	// The following code will generate the Go AST that will simulate a
	// conditional (ternary) operator, in the form of:
	//
	//     noarch.Ternary(
	//         $1,
	//         func () interface{} {
	//             return $2
	//         },
	//         func () interface{} {
	//             return $3
	//         },
	//     )
	//
	// $2 and $3 (the true and false condition respectively) must be wrapped in
	// a closure so that they are not both executed.
	return &goast.CallExpr{
		Fun: goast.NewIdent("noarch.Ternary"),
		Args: []goast.Expr{
			a,
			newTernaryWrapper(b),
			newTernaryWrapper(c),
		},
	}, n.Type, nil
}

// newTernaryWrapper is a helper method used by transpileConditionalOperator().
// It will wrap an expression in a closure.
func newTernaryWrapper(e goast.Expr) *goast.FuncLit {
	return &goast.FuncLit{
		Type: &goast.FuncType{
			Params: &goast.FieldList{},
			Results: &goast.FieldList{
				List: []*goast.Field{
					&goast.Field{
						Type: &goast.InterfaceType{
							Methods: &goast.FieldList{},
						},
					},
				},
			},
		},
		Body: &goast.BlockStmt{
			List: []goast.Stmt{
				&goast.ReturnStmt{
					Results: []goast.Expr{e},
				},
			},
		},
	}
}

// transpileParenExpr transpiles an expression that is wrapped in parentheses.
// There is a special case where "(0)" is treated as a NULL (since that's what
// the macro expands to). We have to return the type as "null" since we don't
// know at this point what the NULL expression will be used in conjuction with.
func transpileParenExpr(n *ast.ParenExpr, p *program.Program) (*goast.ParenExpr, string, error) {
	e, eType, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", err
	}

	r := &goast.ParenExpr{
		X: e,
	}
	if types.IsNullExpr(r) {
		eType = "null"
	}

	return r, eType, nil
}

// getTokenForOperator returns the Go operator token for the provided C
// operator.
func getTokenForOperator(operator string) token.Token {
	switch operator {
	// Arithmetic
	case "--":
		return token.DEC
	case "++":
		return token.INC
	case "+":
		return token.ADD
	case "-":
		return token.SUB
	case "*":
		return token.MUL
	case "/":
		return token.QUO
	case "%":
		return token.REM

	// Assignment
	case "=":
		return token.ASSIGN
	case "+=":
		return token.ADD_ASSIGN
	case "-=":
		return token.SUB_ASSIGN

	// Bitwise
	case "&":
		return token.AND
	case "|":
		return token.OR
	case "~":
		return token.XOR
	case ">>":
		return token.SHR
	case "<<":
		return token.SHL

	// Comparison
	case ">=":
		return token.GEQ
	case "<=":
		return token.LEQ
	case "<":
		return token.LSS
	case ">":
		return token.GTR
	case "!=":
		return token.NEQ
	case "==":
		return token.EQL

	// Logical
	case "!":
		return token.NOT
	case "&&":
		return token.LAND
	case "||":
		return token.LOR
	}

	panic(fmt.Sprintf("unknown operator: %s", operator))
}
