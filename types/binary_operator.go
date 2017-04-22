package types

import (
	"github.com/elliotchance/c2go/program"
)

func ResolveTypeForBinaryOperator(p *program.Program, operator, leftType, rightType string) string {
	if operator == "==" || operator == "!=" {
		return "bool"
	}

	return leftType
}
