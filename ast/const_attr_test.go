package ast

import (
	"testing"
)

func TestConstAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa3b88bbb38 <line:4:1, line:13:1>foo`: &ConstAttr{
			Addr:       0x7fa3b88bbb38,
			Pos:        NewPositionFromString("line:4:1, line:13:1"),
			Tags:       "foo",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
