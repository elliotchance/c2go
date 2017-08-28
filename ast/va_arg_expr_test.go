package ast

import (
	"testing"
)

func TestVAArgExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7ff7d314bca8 <col:6, col:31> 'int *'`: &VAArgExpr{
			Addr:       0x7ff7d314bca8,
			Pos:        NewPositionFromString("col:6, col:31"),
			Type:       "int *",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
