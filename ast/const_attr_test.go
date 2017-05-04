package ast

import (
	"testing"
)

func TestConstAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa3b88bbb38 <line:4:1, line:13:1>foo`: &ConstAttr{
			Address:  "0x7fa3b88bbb38",
			Position: "line:4:1, line:13:1",
			Tags:     "foo",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
