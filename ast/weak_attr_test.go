package ast

import (
	"testing"
)

func TestWeakAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x56069ece5110 <line:736:22>`: &WeakAttr{
			Addr:     0x56069ece5110,
			Position: "line:736:22",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
