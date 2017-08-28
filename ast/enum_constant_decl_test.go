package ast

import (
	"testing"
)

func TestEnumConstantDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x1660db0 <line:185:3> __codecvt_noconv 'int'`: &EnumConstantDecl{
			Addr:       0x1660db0,
			Pos:        NewPositionFromString("line:185:3"),
			Position2:  "",
			Referenced: false,
			Name:       "__codecvt_noconv",
			Type:       "int",
			ChildNodes: []Node{},
		},
		`0x3c77ba8 <line:59:3, col:65> col:3 referenced _ISalnum 'int'`: &EnumConstantDecl{
			Addr:       0x3c77ba8,
			Pos:        NewPositionFromString("line:59:3, col:65"),
			Position2:  "col:3",
			Referenced: true,
			Name:       "_ISalnum",
			Type:       "int",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
