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

func isNullAST(n goast.Expr) bool {
	if p1, ok := n.(*goast.ParenExpr); ok {
		if p2, ok := p1.X.(*goast.BasicLit); ok && p2.Value == "0" {
			return true
		}
	}

	return false
}

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

	// Convert "(0)" to "nil".
	if (operator == token.NEQ || operator == token.EQL) && isNullAST(right) {
		right = goast.NewIdent("nil")
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

		p.AddImport("github.com/elliotchance/c2go/noarch")

		t := types.ResolveType(p, eType)
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

	/*
		if operator == "!" {
		if exprType == "bool" || exprType == "_Bool" {
			return fmt.Sprintf("!(%s)", expr), exprType
		}

		program.AddImport("github.com/elliotchance/c2go/noarch")

		t := types.ResolveType(program, exprType)
		functionName := fmt.Sprintf("noarch.Not%s", util.Ucfirst(t))
		return fmt.Sprintf("%s(%s)", functionName, expr), exprType
	}*/

	return &goast.UnaryExpr{
		X:  e,
		Op: operator,
	}, eType, nil
}

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

func transpileParenExpr(n *ast.ParenExpr, p *program.Program) (*goast.ParenExpr, string, error) {
	e, eType, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", err
	}

	return &goast.ParenExpr{
		Lparen: token.NoPos,
		X:      e,
		Rparen: token.NoPos,
	}, eType, nil
}

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
