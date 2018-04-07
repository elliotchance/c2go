package ast

import (
	"testing"
)

func TestFormatArgAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7f1234567890 <col:47, col:57> 1`: &FormatArgAttr{
			Addr:       0x7f1234567890,
			Pos:        NewPositionFromString("col:47, col:57"),
			Arg:        "1",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
