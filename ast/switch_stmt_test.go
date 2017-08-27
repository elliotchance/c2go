package ast

import (
	"testing"
)

func TestSwitchStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fbca3894638 <line:9:5, line:20:5>`: &SwitchStmt{
			Addr:     0x7fbca3894638,
			Position: "line:9:5, line:20:5",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
