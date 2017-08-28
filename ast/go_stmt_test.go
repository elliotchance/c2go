package ast

import (
	"testing"
)

func TestGotoStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fb9cc1994d8 <line:18893:9, col:14> 'end_getDigits' 0x7fb9cc199490`: &GotoStmt{
			Addr:       0x7fb9cc1994d8,
			Pos:        NewPositionFromString("line:18893:9, col:14"),
			Name:       "end_getDigits",
			Position2:  "0x7fb9cc199490",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
