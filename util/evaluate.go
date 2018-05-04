package util

import (
	goast "go/ast"
	"go/token"
	"strconv"
)

// EvaluateConstExpr evaluates the given expression.
// Returns whether the expr is an integer constant,
// and the resulting number if constant.
func EvaluateConstExpr(expr goast.Expr) (isConst bool, value int64) {
	calc := &calcVisitor{
		isConst: true,
	}
	result := calc.Visit(expr)
	return calc.isConst, result
}

type calcVisitor struct {
	isConst bool
}

func (v *calcVisitor) Visit(node goast.Node) int64 {
	if node == nil {
		return 0
	}
	if be, ok := node.(*goast.BinaryExpr); ok {
		x := v.Visit(be.X)
		y := v.Visit(be.Y)
		switch be.Op {
		case token.ADD:
			return x + y
		case token.SUB:
			return x - y
		case token.MUL:
			return x * y
		case token.QUO:
			if y == 0 {
				v.isConst = false
				return 0
			}
			return x / y
		case token.REM:
			if y == 0 {
				v.isConst = false
				return 0
			}
			return x % y
		case token.AND:
			return x & y
		case token.OR:
			return x | y
		case token.XOR:
			return x ^ y
		case token.SHL:
			return x << uint64(y)
		case token.SHR:
			return x >> uint64(y)
		case token.AND_NOT:
			return x &^ y
		}
	}
	if ue, ok := node.(*goast.UnaryExpr); ok {
		x := v.Visit(ue.X)
		switch ue.Op {
		case token.ADD:
			return x
		case token.SUB:
			return -x
		case token.XOR:
			return ^x
		}
	}
	if ce, ok := node.(*goast.CallExpr); ok {
		if fn, ok2 := ce.Fun.(*goast.Ident); !ok2 || fn.Name != "int32" {
			v.isConst = false
			return 0
		}
		return v.Visit(ce.Args[0])
	}
	if pe, ok := node.(*goast.ParenExpr); ok {
		return v.Visit(pe.X)
	}
	if ie, ok := node.(*goast.BasicLit); ok {
		if ie.Kind == token.INT {
			ret, err := strconv.Atoi(ie.Value)
			if err == nil {
				return int64(ret)
			}
		}
	}
	v.isConst = false
	return 0
}
