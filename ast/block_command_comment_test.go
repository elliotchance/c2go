package ast

import (
	"testing"
)

func TestBlockCommandComment(t *testing.T) {
	nodes := map[string]Node{
		`0x1069fae60 <col:4, line:163:57> Name="abstract"`: &BlockCommandComment{
			Addr:       0x1069fae60,
			Pos:        NewPositionFromString("col:4, line:163:57"),
			Name:       "abstract",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
