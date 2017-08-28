package ast

import (
	"testing"
)

func TestMaxFieldAlignmentAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fd4b7063ac0 <<invalid sloc>> Implicit 32`: &MaxFieldAlignmentAttr{
			Addr:       0x7fd4b7063ac0,
			Pos:        NewPositionFromString("<invalid sloc>"),
			Size:       32,
			ChildNodes: []Node{},
		},
		`0x7fd4b7063ac0 <<invalid sloc>> Implicit 8`: &MaxFieldAlignmentAttr{
			Addr:       0x7fd4b7063ac0,
			Pos:        NewPositionFromString("<invalid sloc>"),
			Size:       8,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
