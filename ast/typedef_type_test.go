package ast

import (
	"testing"
)

func TestTypedefType(t *testing.T) {
	nodes := map[string]Node{
		`0x7f887a0dc760 '__uint16_t' sugar`: &TypedefType{
			Addr:       0x7f887a0dc760,
			Type:       "__uint16_t",
			Tags:       "sugar",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
