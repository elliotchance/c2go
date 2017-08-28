package ast

import (
	"testing"
)

func TestReturnStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fbb7a8325e0 <line:13:4, col:11>`: &ReturnStmt{
			Addr:       0x7fbb7a8325e0,
			Pos:        NewPositionFromString("line:13:4, col:11"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
