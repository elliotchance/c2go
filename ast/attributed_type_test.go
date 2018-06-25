package ast

import (
	"testing"
)

func TestAttributedType(t *testing.T) {
	nodes := map[string]Node{
		`0x2b6c359e30 'int (void) __attribute__((cdecl))' sugar`: &AttributedType{
			Addr:       0x2b6c359e30,
			Type:       "int (void) __attribute__((cdecl))",
			Sugar:      true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
