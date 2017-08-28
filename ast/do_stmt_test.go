package ast

import (
	"testing"
)

func TestDoStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7ff36d0a0938 <line:11:5, line:14:23>`: &DoStmt{
			Addr:       0x7ff36d0a0938,
			Pos:        NewPositionFromString("line:11:5, line:14:23"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
