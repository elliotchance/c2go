package ast

import (
	"testing"
)

func TestHTMLEndTagComment(t *testing.T) {
	nodes := map[string]Node{
		`0x4259670 <col:27, col:30> Name="i"`: &HTMLEndTagComment{
			Addr:       0x4259670,
			Pos:        NewPositionFromString("col:27, col:30"),
			Name:       "i",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
