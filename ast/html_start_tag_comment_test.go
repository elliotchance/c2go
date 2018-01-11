package ast

import (
	"testing"
)

func TestHTMLStartTagComment(t *testing.T) {
	nodes := map[string]Node{
		`0x4259670 <col:27, col:30> Name="i"`: &HTMLStartTagComment{
			Addr:       0x4259670,
			Pos:        NewPositionFromString("col:27, col:30"),
			Name:       "i",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
