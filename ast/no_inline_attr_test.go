package ast

import (
	"testing"
)

func TestNoInlineAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fc02a8a6730 <line:24619:23>`: &NoInlineAttr{
			Addr:       0x7fc02a8a6730,
			Pos:        NewPositionFromString("line:24619:23"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
