package ast

import (
	"testing"
)

func TestDeclStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fb791846e80 <line:11:4, col:31>`: &DeclStmt{
			Addr:       0x7fb791846e80,
			Pos:        NewPositionFromString("line:11:4, col:31"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
