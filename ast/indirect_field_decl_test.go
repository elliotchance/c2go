package ast

import (
	"testing"
)

func TestIndirectFieldDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x2be19a8 <line:167:25> col:25 implicit fpstate 'struct _fpstate *'`: &IndirectFieldDecl{
			Addr:       0x2be19a8,
			Pos:        NewPositionFromString("line:167:25"),
			Position2:  "col:25",
			Implicit:   true,
			Name:       "fpstate",
			Type:       "struct _fpstate *",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
