package ast

import (
	"testing"
)

func TestUnaryOperator(t *testing.T) {
	nodes := map[string]Node{
		`0x7fe0260f50d8 <col:6, col:12> 'int' prefix '--'`: &UnaryOperator{
			Addr:       0x7fe0260f50d8,
			Pos:        NewPositionFromString("col:6, col:12"),
			Type:       "int",
			IsLvalue:   false,
			IsPrefix:   true,
			Operator:   "--",
			ChildNodes: []Node{},
		},
		`0x7fe0260fb468 <col:11, col:18> 'unsigned char' lvalue prefix '*'`: &UnaryOperator{
			Addr:       0x7fe0260fb468,
			Pos:        NewPositionFromString("col:11, col:18"),
			Type:       "unsigned char",
			IsLvalue:   true,
			IsPrefix:   true,
			Operator:   "*",
			ChildNodes: []Node{},
		},
		`0x7fe0260fb448 <col:12, col:18> 'unsigned char *' postfix '++'`: &UnaryOperator{
			Addr:       0x7fe0260fb448,
			Pos:        NewPositionFromString("col:12, col:18"),
			Type:       "unsigned char *",
			IsLvalue:   false,
			IsPrefix:   false,
			Operator:   "++",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
