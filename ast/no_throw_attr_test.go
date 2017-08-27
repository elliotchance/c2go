package ast

import (
	"testing"
)

func TestNoThrowAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa1488273a0 <line:7:4, line:11:4>`: &NoThrowAttr{
			Addr:       0x7fa1488273a0,
			Pos:        "line:7:4, line:11:4",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
