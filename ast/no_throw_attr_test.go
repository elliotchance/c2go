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
			Inherited:  false,
			Implicit:   false,
		},
		`0x5605ceaf4b88 <col:12> Implicit`: &NoThrowAttr{
			Addr:       0x5605ceaf4b88,
			Pos:        NewPositionFromString("col:12"),
			ChildNodes: []Node{},
			Inherited:  false,
			Implicit:   true,
		},
		`0x4153c50 </usr/include/unistd.h:779:46> Inherited`: &NoThrowAttr{
			Addr:       0x4153c50,
			Pos:        NewPositionFromString("/usr/include/unistd.h:779:46"),
			ChildNodes: []Node{},
			Inherited:  true,
			Implicit:   false,
		},
		`0x1038b8828 <col:20> Inherited Implicit`: &NoThrowAttr{
			Addr:       0x1038b8828,
			Pos:        NewPositionFromString("col:20"),
			ChildNodes: []Node{},
			Inherited:  true,
			Implicit:   true,
		},
	}

	runNodeTests(t, nodes)
}
