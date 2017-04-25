package ast

import (
	"testing"
)

func TestIntegerLiteral(t *testing.T) {
	nodes := map[string]Node{
		`0x7fbe9804bcc8 <col:14> 'int' 1`: &IntegerLiteral{
			Address:  "0x7fbe9804bcc8",
			Position: "col:14",
			Type:     "int",
			Value:    1,
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
