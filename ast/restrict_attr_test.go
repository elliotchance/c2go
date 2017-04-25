package ast

import (
	"testing"
)

func TestRestrictAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7f980b858305 <line:11:7, line:18:7> foo`: &RestrictAttr{
			Address:  "0x7f980b858305",
			Position: "line:11:7, line:18:7",
			Name:     "foo",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
