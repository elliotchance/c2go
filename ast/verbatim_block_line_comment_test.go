package ast

import (
	"testing"
)

func TestVerbatimBlockLineComment(t *testing.T) {
	nodes := map[string]Node{
		`0x10f8e8e20 <col:39, col:54> Text=" OSAtomicAdd32}"`: &VerbatimBlockLineComment{
			Addr:       0x10f8e8e20,
			Pos:        NewPositionFromString("col:39, col:54"),
			Text:       " OSAtomicAdd32}",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
