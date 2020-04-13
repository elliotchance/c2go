package ast

import (
	"testing"
)

func TestEnumDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x22a6c80 <line:180:1, line:186:1> __codecvt_result`: &EnumDecl{
			Addr:       0x22a6c80,
			Pos:        NewPositionFromString("line:180:1, line:186:1"),
			Position2:  "",
			Name:       "__codecvt_result",
			ChildNodes: []Node{},
		},
		`0x32fb5a0 <enum.c:3:1, col:45> col:6 week`: &EnumDecl{
			Addr:       0x32fb5a0,
			Pos:        NewPositionFromString("enum.c:3:1, col:45"),
			Position2:  " col:6",
			Name:       "week",
			ChildNodes: []Node{},
		},
		`0x5570ff477388 prev 0x5570ff445c40 <line:1030:1, line:1034:1> line:1030:6 cb_assign_type`: &EnumDecl{
			Addr:       0x5570ff477388,
			Prev:       0x5570ff445c40,
			Pos:        NewPositionFromString("line:1030:1, line:1034:1"),
			Position2:  " line:1030:6",
			Name:       "cb_assign_type",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
