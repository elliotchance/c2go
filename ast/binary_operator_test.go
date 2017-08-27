package ast

import (
	"testing"
)

func TestBinaryOperator(t *testing.T) {
	nodes := map[string]Node{
		`0x7fca2d8070e0 <col:11, col:23> 'unsigned char' '='`: &BinaryOperator{
			Addr:       0x7fca2d8070e0,
			Position:   "col:11, col:23",
			Type:       "unsigned char",
			Operator:   "=",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
