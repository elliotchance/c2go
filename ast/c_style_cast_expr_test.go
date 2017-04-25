package ast

import (
	"testing"
)

func TestCStyleCastExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fddc18fb2e0 <col:50, col:56> 'char' <IntegralCast>`: &CStyleCastExpr{
			Address:  "0x7fddc18fb2e0",
			Position: "col:50, col:56",
			Type:     "char",
			Kind:     "IntegralCast",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
