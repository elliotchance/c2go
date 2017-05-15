package ast

import (
	"testing"
)

func TestAlignedAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7f8a1d8ccfd0 <col:47, col:57> aligned`: &AlignedAttr{
			Address:  "0x7f8a1d8ccfd0",
			Position: "col:47, col:57",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
