package ast

import (
	"testing"
)

func TestArraySubscriptExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fe35b85d180 <col:63, col:69> 'char *' lvalue`: &ArraySubscriptExpr{
			Addr:       0x7fe35b85d180,
			Pos:        NewPositionFromString("col:63, col:69"),
			Type:       "char *",
			Kind:       "lvalue",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
