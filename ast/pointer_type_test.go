package ast

import (
	"testing"
)

func TestPointerType(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa3b88bbb30 'struct _opaque_pthread_t *'`: &PointerType{
			Addr:       0x7fa3b88bbb30,
			Type:       "struct _opaque_pthread_t *",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
