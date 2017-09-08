package ast

import (
	"testing"
)

func TestGCCAsmStmtStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fad830c9e38 <line:13:5, col:57>`: &GCCAsmStmt{
			Addr:       0x7fad830c9e38,
			Pos:        NewPositionFromString("line:13:5, col:57"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
