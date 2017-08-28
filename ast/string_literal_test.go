package ast

import (
	"testing"
)

func TestStringLiteral(t *testing.T) {
	nodes := map[string]Node{
		`0x7fe16f0b4d58 <col:11> 'char [45]' lvalue "Number of command line arguments passed: %d\n"`: &StringLiteral{
			Addr:       0x7fe16f0b4d58,
			Pos:        NewPositionFromString("col:11"),
			Type:       "char [45]",
			Lvalue:     true,
			Value:      "Number of command line arguments passed: %d\n",
			ChildNodes: []Node{},
		},
		`0x22ac548 <col:14> 'char [14]' lvalue "x\vx\000xxx\axx\tx\n"`: &StringLiteral{
			Addr:       0x22ac548,
			Pos:        NewPositionFromString("col:14"),
			Type:       "char [14]",
			Lvalue:     true,
			Value:      "x\vx\x00xxx\axx\tx\n",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
