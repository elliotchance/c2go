package ast

import (
	"testing"
)

func TestForStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7f961e018848 <line:9:4, line:10:70>`: &ForStmt{
			Addr:       0x7f961e018848,
			Pos:        NewPositionFromString("line:9:4, line:10:70"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
