package transpiler

import (
	"fmt"
	"go/token"
)

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
	case "%":
		return token.REM
	case "*":
		return token.MUL

	// Assignment
	case "=":
		return token.ASSIGN

	// Bitwise
	case "&":
		return token.AND

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
	case "&&":
		return token.LAND
	case "||":
		return token.LOR
	}

	panic(fmt.Sprintf("unknown operator: %s", operator))
}
