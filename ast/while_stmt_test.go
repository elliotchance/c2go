package ast

import (
	"testing"
)

func TestWhileStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa1478273a0 <line:7:4, line:11:4>`: &WhileStmt{
			Addr:     0x7fa1478273a0,
			Position: "line:7:4, line:11:4",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
