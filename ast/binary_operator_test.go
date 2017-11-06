package ast

import (
	"testing"
)

func TestBinaryOperator(t *testing.T) {
	nodes := map[string]Node{
		`0x7fca2d8070e0 <col:11, col:23> 'unsigned char' '='`: &BinaryOperator{
			Addr:       0x7fca2d8070e0,
			Pos:        NewPositionFromString("col:11, col:23"),
			Type:       "unsigned char",
			Operator:   "=",
			ChildNodes: []Node{},
		},
		`0x1ff95b8 <line:78:2, col:7> 'T_ENUM':'T_ENUM' '='`: &BinaryOperator{
			Addr:       0x1ff95b8,
			Pos:        NewPositionFromString("line:78:2, col:7"),
			Type:       "T_ENUM",
			Type2:      "T_ENUM",
			Operator:   "=",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
