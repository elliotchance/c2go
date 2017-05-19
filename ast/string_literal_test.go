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
		`0x22ac548 <col:14> 'char [14]' lvalue "x\vx\000xxx\axx\tx\n"`: &StringLiteral{
			Address:  "0x22ac548",
			Position: "col:14",
			Type:     "char [14]",
			Lvalue:   true,
			Value:    "x\vx\x00xxx\axx\tx\n",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
