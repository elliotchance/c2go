package ast

import (
	"testing"
)

func TestConditionalOperator(t *testing.T) {
	nodes := map[string]Node{
		`0x7fc6ae0bc678 <col:6, col:89> 'void'`: &ConditionalOperator{
			Address:  "0x7fc6ae0bc678",
			Position: "col:6, col:89",
			Type:     "void",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
