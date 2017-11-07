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
			ChildNodes: []Node{},
		},
		`0x1ff8708 <col:14, col:17> 'T_ENUM':'T_ENUM' lvalue`: &ParenExpr{
			Addr:       0x1ff8708,
			Pos:        NewPositionFromString("col:14, col:17"),
			Type:       "T_ENUM",
			Type2:      "T_ENUM",
			Lvalue:     true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
