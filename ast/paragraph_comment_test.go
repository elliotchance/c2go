package ast

import (
	"testing"
)

func TestParagraphComment(t *testing.T) {
	nodes := map[string]Node{
		`0x3860920 <line:10176:4, line:10180:45>`: &ParagraphComment{
			Addr:       0x3860920,
			Pos:        NewPositionFromString("line:10176:4, line:10180:45"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
