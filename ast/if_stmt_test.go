package ast

import (
	"testing"
)

func TestIfStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fc0a69091d0 <line:11:7, line:18:7>`: &IfStmt{
			Addr:       0x7fc0a69091d0,
			Pos:        NewPositionFromString("line:11:7, line:18:7"),
			ChildNodes: []Node{},
		},
		`0x560d43f4ee98 <line:65:2, line:72:2> has_else`: &IfStmt{
			Addr:       0x560d43f4ee98,
			Pos:        NewPositionFromString("line:65:2, line:72:2"),
			HasElse:    true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
