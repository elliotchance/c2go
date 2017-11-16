package ast

import (
	"testing"
)

func TestNoThrowAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa1488273a0 <line:7:4, line:11:4>`: &NoThrowAttr{
			Addr:       0x7fa1488273a0,
			Pos:        NewPositionFromString("line:7:4, line:11:4"),
			ChildNodes: []Node{},
			Implicit:   false,
		},
		`0x5605ceaf4b88 <col:12> Implicit`: &NoThrowAttr{
			Addr:       0x5605ceaf4b88,
			Pos:        NewPositionFromString("col:12"),
			ChildNodes: []Node{},
			Implicit:   true,
		},
	}

	runNodeTests(t, nodes)
}
