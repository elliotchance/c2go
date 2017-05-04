package ast

import (
	"testing"
)

func TestStringLiteral(t *testing.T) {
	nodes := map[string]Node{
		`0x7fe16f0b4d58 <col:11> 'char [45]' lvalue "Number of command line arguments passed: %d\n"`: &StringLiteral{
			Address:  "0x7fe16f0b4d58",
			Position: "col:11",
			Type:     "char [45]",
			Lvalue:   true,
			Value:    "Number of command line arguments passed: %d\n",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
