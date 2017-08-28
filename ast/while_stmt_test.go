package ast

import (
	"testing"
)

func TestWhileStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa1478273a0 <line:7:4, line:11:4>`: &WhileStmt{
			Addr:       0x7fa1478273a0,
			Pos:        NewPositionFromString("line:7:4, line:11:4"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
