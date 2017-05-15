package ast

import (
	"testing"
)

func TestNoInlineAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fc02a8a6730 <line:24619:23>`: &NoInlineAttr{
			Address:  "0x7fc02a8a6730",
			Position: "line:24619:23",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
