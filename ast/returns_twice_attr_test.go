package ast

import (
	"testing"
)

func TestReturnsTwiceAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7ff8e9091640 <col:7> Implicit`: &ReturnsTwiceAttr{
			Addr:       0x7ff8e9091640,
			Pos:        NewPositionFromString("col:7"),
			ChildNodes: []Node{},
			Inherited:  false,
			Implicit:   true,
		},
		`0x564a73a5ccc8 <col:16> Inherited Implicit`: &ReturnsTwiceAttr{
			Addr:       0x564a73a5ccc8,
			Pos:        NewPositionFromString("col:16"),
			ChildNodes: []Node{},
			Inherited:  true,
			Implicit:   true,
		},
	}

	runNodeTests(t, nodes)
}
