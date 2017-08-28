package ast

import (
	"testing"
)

func TestBreakStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fca2d8070e0 <col:11, col:23>`: &BreakStmt{
			Addr:       0x7fca2d8070e0,
			Pos:        NewPositionFromString("col:11, col:23"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
