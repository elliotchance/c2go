package ast

import (
	"testing"
)

func TestInitListExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fbdd1906c20 <col:52, line:17160:1> 'const unsigned char [256]'`: &InitListExpr{
			Addr:       0x7fbdd1906c20,
			Pos:        NewPositionFromString("col:52, line:17160:1"),
			Type:       "const unsigned char [256]",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
