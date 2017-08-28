package ast

import (
	"testing"
)

func TestContinueStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x1e044e0 <col:20>`: &ContinueStmt{
			Addr:       0x1e044e0,
			Pos:        NewPositionFromString("col:20"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
