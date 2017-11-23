package ast

import (
	"testing"
)

func TestConstantArrayType(t *testing.T) {
	nodes := map[string]Node{
		`0x7f94ad016a40 'struct __va_list_tag [1]' 1 `: &ConstantArrayType{
			Addr:       0x7f94ad016a40,
			Type:       "struct __va_list_tag [1]",
			Size:       1,
			ChildNodes: []Node{},
		},
		`0x7f8c5f059d20 'char [37]' 37 `: &ConstantArrayType{
			Addr:       0x7f8c5f059d20,
			Type:       "char [37]",
			Size:       37,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
