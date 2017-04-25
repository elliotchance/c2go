package ast

import (
	"testing"
)

func TestEnumConstantDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x1660db0 <line:185:3> __codecvt_noconv 'int'`: &EnumConstantDecl{
			Address:   "0x1660db0",
			Position:  "line:185:3",
			Position2: "",
			Name:      "__codecvt_noconv",
			Type:      "int",
			Children:  []Node{},
		},
	}

	runNodeTests(t, nodes)
}
