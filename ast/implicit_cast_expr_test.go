package ast

import (
	"testing"
)

func TestImplicitCastExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7f9f5b0a1288 <col:8> 'FILE *' <LValueToRValue>`: &ImplicitCastExpr{
			Addr:     0x7f9f5b0a1288,
			Position: "col:8",
			Type:     "FILE *",
			Kind:     "LValueToRValue",
			Children: []Node{},
		},
		`0x7f9f5b0a7828 <col:11> 'int (*)(int, FILE *)' <FunctionToPointerDecay>`: &ImplicitCastExpr{
			Addr:     0x7f9f5b0a7828,
			Position: "col:11",
			Type:     "int (*)(int, FILE *)",
			Kind:     "FunctionToPointerDecay",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
