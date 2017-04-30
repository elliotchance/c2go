package transpiler

import (
	"fmt"
	"go/token"

	goast "go/ast"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
)

func transpileBinaryOperator(n *ast.BinaryOperator, p *program.Program) (*goast.BinaryExpr, string, error) {
	left, _, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", err
	}

	right, _, err := transpileToExpr(n.Children[1], p)
	if err != nil {
		return nil, "", err
	}

	return &goast.BinaryExpr{
		X:     left,
		OpPos: token.NoPos,
		Op:    getTokenForOperator(n.Operator),
		Y:     right,
	}, "", nil
}

func transpileUnaryOperator(n *ast.UnaryOperator, p *program.Program) (*goast.UnaryExpr, string, error) {
	left, _, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", err
	}

	return &goast.UnaryExpr{
		X:     left,
		OpPos: token.NoPos,
		Op:    getTokenForOperator(n.Operator),
	}, "", nil
}

func transpileConditionalOperator(n *ast.ConditionalOperator, p *program.Program) (*goast.CallExpr, string, error) {
	// TODO: check errors for these
	a, _, _ := transpileToExpr(n.Children[0], p)
	b, _, _ := transpileToExpr(n.Children[1], p)
	c, _, _ := transpileToExpr(n.Children[2], p)

	p.AddImport("github.com/elliotchance/c2go/noarch")

	return &goast.CallExpr{
		Fun:      goast.NewIdent("noarch.Ternary"),
		Lparen:   token.NoPos,
		Args:     []goast.Expr{a, b, c},
		Ellipsis: token.NoPos,
		Rparen:   token.NoPos,
	}, "", nil

	// src := fmt.Sprintf("noarch.Ternary(%s, func () interface{} { return %s }, func () interface{} { return %s })", a, b, c)
	// return src, n.Type
}

func transpileParenExpr(n *ast.ParenExpr, p *program.Program) (*goast.ParenExpr, string, error) {
	e, _, err := transpileToExpr(n.Children[0], p)
	if err != nil {
		return nil, "", err
	}

	return &goast.ParenExpr{
		Lparen: token.NoPos,
		X:      e,
		Rparen: token.NoPos,
	}, "", nil
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

	// Bitwise
	case "&":
		return token.AND
	case "|":
		return token.OR

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
