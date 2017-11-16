package ast

import (
	"testing"
)

func TestVisibilityAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x55c49d8dd1d8 <col:16, col:36> Default`: &VisibilityAttr{
			Addr:       0x55c49d8dd1d8,
			Pos:        NewPositionFromString("col:16, col:36"),
			ChildNodes: []Node{},
			IsDefault:  true,
		},
	}

	runNodeTests(t, nodes)
}
