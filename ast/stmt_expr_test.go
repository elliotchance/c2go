package ast

import (
	"testing"
)

func TestStmtExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7ff4f9100d28 <col:11, col:18> 'int'`: &StmtExpr{
			Addr:       0x7ff4f9100d28,
			Pos:        NewPositionFromString("col:11, col:18"),
			Type:       "int",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
