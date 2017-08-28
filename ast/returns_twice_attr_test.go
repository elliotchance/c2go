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
		},
	}

	runNodeTests(t, nodes)
}
