package ast

import (
	"testing"
)

func TestTextComment(t *testing.T) {
	nodes := map[string]Node{
		`0x3085bc0 <line:9950:2, col:29> Text="* CUSTOM AUXILIARY FUNCTIONS"`: &TextComment{
			Addr:       0x3085bc0,
			Pos:        NewPositionFromString("line:9950:2, col:29"),
			Text:       "* CUSTOM AUXILIARY FUNCTIONS",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
