package types

import (
	"github.com/elliotchance/c2go/program"
)

func ResolveTypeForBinaryOperator(p *program.Program, operator, leftType, rightType string) string {
	if operator == "==" ||
		operator == "!=" ||
		operator == ">" ||
		operator == ">=" ||
		operator == "<" ||
		operator == "<=" {
		return "bool"
	}

	/*
		Example of c code:
		int fact(int n) {
			int lcv, p;
			for(p=1, lcv=2; lcv <= n; p=p*lcv, lcv++); // <<== Here 2 timer in initialization and increment
			return p;
		}
	*/
	if operator == "," {
		panic("Now, Algorithm is not fully ready")
	}

	return leftType
}
