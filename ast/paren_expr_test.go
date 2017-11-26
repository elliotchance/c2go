package ast

import (
	"testing"
)

func TestParenExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fb0bc8b2308 <col:10, col:25> 'unsigned char'`: &ParenExpr{
			Addr:       0x7fb0bc8b2308,
			Pos:        NewPositionFromString("col:10, col:25"),
			Type:       "unsigned char",
			Type2:      "",
			Lvalue:     false,
			IsBitfield: false,
			ChildNodes: []Node{},
		},
		`0x1ff8708 <col:14, col:17> 'T_ENUM':'T_ENUM' lvalue`: &ParenExpr{
			Addr:       0x1ff8708,
			Pos:        NewPositionFromString("col:14, col:17"),
			Type:       "T_ENUM",
			Type2:      "T_ENUM",
			Lvalue:     true,
			IsBitfield: false,
			ChildNodes: []Node{},
		},
		`0x55efc60798b0 <col:15, col:27> 'bft':'unsigned int' lvalue bitfield`: &ParenExpr{
			Addr:       0x55efc60798b0,
			Pos:        NewPositionFromString("col:15, col:27"),
			Type:       "bft",
			Type2:      "unsigned int",
			Lvalue:     true,
			IsBitfield: true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
