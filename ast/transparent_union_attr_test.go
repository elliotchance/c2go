package ast

import (
	"testing"
)

func TestTransparentUnionAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x304f700 <col:35>`: &TransparentUnionAttr{
			Addr:       0x304f700,
			Pos:        NewPositionFromString("col:35"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
