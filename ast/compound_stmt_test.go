package ast

import (
	"testing"
)

func TestCompoundStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fbd0f014f18 <col:54, line:358:1>`: &CompoundStmt{
			Addr:       0x7fbd0f014f18,
			Pos:        NewPositionFromString("col:54, line:358:1"),
			ChildNodes: []Node{},
		},
		`0x7fbd0f8360b8 <line:4:1, line:13:1>`: &CompoundStmt{
			Addr:       0x7fbd0f8360b8,
			Pos:        NewPositionFromString("line:4:1, line:13:1"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
