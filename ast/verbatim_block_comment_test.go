package ast

import (
	"testing"
)

func TestVerbatimBlockComment(t *testing.T) {
	nodes := map[string]Node{
		`0x107781dd0 <col:34, col:39> Name="link" CloseName=""`: &VerbatimBlockComment{
			Addr:       0x107781dd0,
			Pos:        NewPositionFromString("col:34, col:39"),
			Name:       "link",
			CloseName:  "",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
