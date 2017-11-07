package ast

import (
	"testing"
)

func TestCStyleCastExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fddc18fb2e0 <col:50, col:56> 'char' <IntegralCast>`: &CStyleCastExpr{
			Addr:       0x7fddc18fb2e0,
			Pos:        NewPositionFromString("col:50, col:56"),
			Type:       "char",
			Kind:       "IntegralCast",
			ChildNodes: []Node{},
		},
		`0x2781518 <col:7, col:17> 'T_ENUM':'T_ENUM' <IntegralCast>`: &CStyleCastExpr{
			Addr:       0x2781518,
			Pos:        NewPositionFromString("col:7, col:17"),
			Type:       "T_ENUM",
			Type2:      "T_ENUM",
			Kind:       "IntegralCast",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
