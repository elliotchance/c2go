package ast

import (
	"testing"
)

func TestConstantExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x5558ffde68f0 <col:34, col:55> 'unsigned long'`: &ConstantExpr{
			Addr:       0x5558ffde68f0,
			Pos:        NewPositionFromString("col:34, col:55"),
			Type:       "unsigned long",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
