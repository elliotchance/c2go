package ast

import (
	"testing"
)

func TestConstantExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7faf4904eed8 <col:31> 'int'`: &ConstantExpr{
			Addr:       0x7faf4904eed8,
			Pos:        NewPositionFromString("col:31"),
			Type:       "int",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
