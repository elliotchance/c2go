package ast

import (
	"testing"
)

func TestNoEscapeAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7febb582f4c0 <col:75>`: &NoEscapeAttr{
			Addr:       0x7febb582f4c0,
			Pos:        NewPositionFromString("col:75"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
