package ast

import (
	"testing"
)

func TestMemberExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fcc758e34a0 <col:8, col:12> 'int' lvalue ->_w 0x7fcc758d60c8`: &MemberExpr{
			Address:   "0x7fcc758e34a0",
			Position:  "col:8, col:12",
			Type:      "int",
			Lvalue:    true,
			Name:      "_w",
			Address2:  "0x7fcc758d60c8",
			IsPointer: true,
			Children:  []Node{},
		},
		`0x7fcc76004210 <col:12, col:16> 'unsigned char *' lvalue ->_p 0x7fcc758d6018`: &MemberExpr{
			Address:   "0x7fcc76004210",
			Position:  "col:12, col:16",
			Type:      "unsigned char *",
			Lvalue:    true,
			Name:      "_p",
			Address2:  "0x7fcc758d6018",
			IsPointer: true,
			Children:  []Node{},
		},
		`0x7f85338325b0 <col:4, col:13> 'float' lvalue .constant 0x7f8533832260`: &MemberExpr{
			Address:   "0x7f85338325b0",
			Position:  "col:4, col:13",
			Type:      "float",
			Lvalue:    true,
			Name:      "constant",
			Address2:  "0x7f8533832260",
			IsPointer: false,
			Children:  []Node{},
		},
		`0x7f8533832670 <col:4, col:13> 'char *' lvalue .pointer 0x7f85338322b8`: &MemberExpr{
			Address:   "0x7f8533832670",
			Position:  "col:4, col:13",
			Type:      "char *",
			Lvalue:    true,
			Name:      "pointer",
			Address2:  "0x7f85338322b8",
			IsPointer: false,
			Children:  []Node{},
		},
	}

	runNodeTests(t, nodes)
}
