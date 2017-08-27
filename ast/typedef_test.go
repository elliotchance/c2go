package ast

import (
	"testing"
)

func TestTypedef(t *testing.T) {
	nodes := map[string]Node{
		`0x7f84d10dc1d0 '__darwin_ssize_t'`: &Typedef{
			Addr:       0x7f84d10dc1d0,
			Type:       "__darwin_ssize_t",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
