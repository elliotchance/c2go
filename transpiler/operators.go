package transpiler

import (
	"fmt"
	"go/token"
)

func getTokenForOperator(operator string) token.Token {
	switch operator {
	case "&":
		return token.AND
	case "--":
		return token.DEC
	case "++":
		return token.INC
	case "*":
		return token.MUL
	case ">=":
		return token.GEQ
	case "<=":
		return token.LEQ
	case "!=":
		return token.NEQ
	case "&&":
		return token.LAND
	case "||":
		return token.LOR
	case "=":
		return token.ASSIGN
	case "==":
		return token.EQL
	case "-":
		return token.SUB
	case "%":
		return token.REM
	}

	panic(fmt.Sprintf("unknown operator: %s", operator))
}
