package ast

import (
	"testing"
)

func TestEnum(t *testing.T) {
	nodes := map[string]Node{
		`0x7f980b858308 'foo'`: &Enum{
			Addr:       0x7f980b858308,
			Name:       "foo",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
