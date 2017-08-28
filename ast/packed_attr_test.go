package ast

import (
	"testing"
)

func TestPackedAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fae33b1ed40 <line:551:18>`: &PackedAttr{
			Addr:       0x7fae33b1ed40,
			Pos:        NewPositionFromString("line:551:18"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
