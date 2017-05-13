package ast

import (
	"testing"
)

func TestUnaryExprOrTypeTraitExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fccd70adf50 <col:29, col:40> 'unsigned long' sizeof 'char'`: &UnaryExprOrTypeTraitExpr{
			Address:  "0x7fccd70adf50",
			Position: "col:29, col:40",
			Type1:    "unsigned long",
			Function: "sizeof",
			Type2:    "char",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
