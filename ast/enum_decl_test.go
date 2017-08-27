package ast

import (
	"testing"
)

func TestEnumDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x22a6c80 <line:180:1, line:186:1> __codecvt_result`: &EnumDecl{
			Addr:       0x22a6c80,
			Position:   "line:180:1, line:186:1",
			Position2:  "",
			Name:       "__codecvt_result",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
