package ast

import (
	"testing"
)

func TestUnusedAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fe3e01416d0 <col:47> unused`: &UnusedAttr{
			Addr:       0x7fe3e01416d0,
			Pos:        NewPositionFromString("col:47"),
			ChildNodes: []Node{},
			IsUnused:   true,
		},
		`0x7fe3e01416d0 <col:47>`: &UnusedAttr{
			Addr:       0x7fe3e01416d0,
			Pos:        NewPositionFromString("col:47"),
			ChildNodes: []Node{},
			IsUnused:   false,
		},
	}

	runNodeTests(t, nodes)
}
