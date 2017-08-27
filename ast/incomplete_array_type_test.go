package ast

import (
	"testing"
)

func TestIncompleteArrayType(t *testing.T) {
	nodes := map[string]Node{
		`0x7fcb7d005c20 'int []' `: &IncompleteArrayType{
			Addr:       0x7fcb7d005c20,
			Type:       "int []",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
