package ast

import (
	"testing"
)

func TestNotTailCalledAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fc8fa094558 <col:107> `: &NotTailCalledAttr{
			Addr:       0x7fc8fa094558,
			Pos:        NewPositionFromString("col:107"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
