package ast

import (
	"testing"
)

func TestFunctionNoProtoType(t *testing.T) {
	nodes := map[string]Node{
		`0x556e32bfde50 'int ()' cdecl`: &FunctionNoProtoType{
			Addr:        0x556e32bfde50,
			Type:        "int ()",
			CallingConv: "cdecl",
			ChildNodes:  []Node{},
		},
	}

	runNodeTests(t, nodes)
}
