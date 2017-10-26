package ast

import (
	"testing"
)

func TestAllocSizeAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7f8e390a5d38 <col:100, col:114> 1 2`: &AllocSizeAttr{
			Addr:       0x7f8e390a5d38,
			Pos:        NewPositionFromString("col:100, col:114"),
			A:          1,
			B:          2,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
