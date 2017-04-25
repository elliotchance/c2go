package ast

import (
	"testing"
)

func TestCompoundStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fbd0f014f18 <col:54, line:358:1>`: &CompoundStmt{
			Address:  "0x7fbd0f014f18",
			Position: "col:54, line:358:1",
			Children: []Node{},
		},
		`0x7fbd0f8360b8 <line:4:1, line:13:1>`: &CompoundStmt{
			Address:  "0x7fbd0f8360b8",
			Position: "line:4:1, line:13:1",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
