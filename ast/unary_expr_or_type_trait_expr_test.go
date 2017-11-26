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
			Type3:      "",
			ChildNodes: []Node{},
		},
		`0x7fae1a800190 <col:36, col:44> 'unsigned long' sizeof`: &UnaryExprOrTypeTraitExpr{
			Addr:       0x7fae1a800190,
			Pos:        NewPositionFromString("col:36, col:44"),
			Type1:      "unsigned long",
			Function:   "sizeof",
			Type2:      "",
			Type3:      "",
			ChildNodes: []Node{},
		},
		`0x557e575e70b8 <col:432, col:452> 'unsigned long' sizeof 'union MyUnion':'union MyUnion'`: &UnaryExprOrTypeTraitExpr{
			Addr:       0x557e575e70b8,
			Pos:        NewPositionFromString("col:432, col:452"),
			Type1:      "unsigned long",
			Function:   "sizeof",
			Type2:      "union MyUnion",
			Type3:      "union MyUnion",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
