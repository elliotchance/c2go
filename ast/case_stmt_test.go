package ast

import (
	"testing"
)

func TestCaseStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fc8b5094688 <line:11:5, line:12:21>`: &CaseStmt{
			Addr:     0x7fc8b5094688,
			Position: "line:11:5, line:12:21",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
