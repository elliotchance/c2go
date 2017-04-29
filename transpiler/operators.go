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
