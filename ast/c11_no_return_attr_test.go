package ast

import (
	"testing"
)

func TestC11NoReturnAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x55a5fc736cf0 <col:1>`: &C11NoReturnAttr{
			Addr:       0x55a5fc736cf0,
			Pos:        NewPositionFromString("col:1"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
