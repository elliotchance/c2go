package ast

import (
	"testing"
)

func Test(t *testing.T) {
	nodes := map[string]Node{
		`0x7f980b858309 <line:11:7, line:18:7> foo`: &ModeAttr{
			Addr:     0x7f980b858309,
			Position: "line:11:7, line:18:7",
			Name:     "foo",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
