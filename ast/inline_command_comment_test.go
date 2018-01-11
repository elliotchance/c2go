package ast

import (
	"testing"
)

func TestInlineCommandComment(t *testing.T) {
	nodes := map[string]Node{
		`0x22e3510 <col:2, col:6> Name="NOTE" RenderNormal`: &InlineCommandComment{
			Addr:       0x22e3510,
			Pos:        NewPositionFromString("col:2, col:6"),
			Other:      "Name=\"NOTE\" RenderNormal",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
