package ast

import (
	"testing"
)

func TestUnaryExprOrTypeTraitExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fccd70adf50 <col:29, col:40> 'unsigned long' sizeof 'char'`: &UnaryExprOrTypeTraitExpr{
			Addr:       0x7fccd70adf50,
			Pos:        NewPositionFromString("col:29, col:40"),
			Type1:      "unsigned long",
			Function:   "sizeof",
			Type2:      "char",
			ChildNodes: []Node{},
		},
		`0x7fae1a800190 <col:36, col:44> 'unsigned long' sizeof`: &UnaryExprOrTypeTraitExpr{
			Addr:       0x7fae1a800190,
			Pos:        NewPositionFromString("col:36, col:44"),
			Type1:      "unsigned long",
			Function:   "sizeof",
			Type2:      "",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
