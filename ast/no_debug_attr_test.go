package ast

import (
	"testing"
)

func TestNoDebugAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x33b8788 <col:59>`: &NoDebugAttr{
			Addr:       0x33b8788,
			Pos:        NewPositionFromString("col:59"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
