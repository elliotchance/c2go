package ast

import (
	"testing"
)

func TestEnumConstantDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x1660db0 <line:185:3> __codecvt_noconv 'int'`: &EnumConstantDecl{
			Address:    "0x1660db0",
			Position:   "line:185:3",
			Position2:  "",
			Referenced: false,
			Name:       "__codecvt_noconv",
			Type:       "int",
			Children:   []Node{},
		},
		`0x3c77ba8 <line:59:3, col:65> col:3 referenced _ISalnum 'int'`: &EnumConstantDecl{
			Address:    "0x3c77ba8",
			Position:   "line:59:3, col:65",
			Position2:  "col:3",
			Referenced: true,
			Name:       "_ISalnum",
			Type:       "int",
			Children:   []Node{},
		},
	}

	runNodeTests(t, nodes)
}
