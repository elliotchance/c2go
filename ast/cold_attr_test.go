package ast

import (
	"testing"
)

func TestColdAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fbdda09a3b8 <col:42>`: &ColdAttr{
			Addr:       0x7fbdda09a3b8,
			Pos:        NewPositionFromString("col:42"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
