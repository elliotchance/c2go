package ast

import (
	"testing"
)

func TestUnusedAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fe3e01416d0 <col:47> unused`: &UnusedAttr{
			Addr:       0x7fe3e01416d0,
			Position:   "col:47",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
