package ast

import (
	"testing"
)

func TestIntegerLiteral(t *testing.T) {
	nodes := map[string]Node{
		`0x7fbe9804bcc8 <col:14> 'int' 1`: &IntegerLiteral{
			Addr:       0x7fbe9804bcc8,
			Pos:        NewPositionFromString("col:14"),
			Type:       "int",
			Value:      "1",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
