package ast

import (
	"testing"
)

func TestDoStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7ff36d0a0938 <line:11:5, line:14:23>`: &DoStmt{
			Address:  "0x7ff36d0a0938",
			Position: "line:11:5, line:14:23",
			Children: []interface{}{},
		},
	}

	runNodeTests(t, nodes)
}
