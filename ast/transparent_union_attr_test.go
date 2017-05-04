package ast

import (
	"testing"
)

func TestTransparentUnionAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x304f700 <col:35>`: &TransparentUnionAttr{
			Address:  "0x304f700",
			Position: "col:35",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
