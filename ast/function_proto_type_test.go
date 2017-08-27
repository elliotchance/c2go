package ast

import (
	"testing"
)

func TestFunctionProtoType(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa3b88bbb30 'struct _opaque_pthread_t *' foo`: &FunctionProtoType{
			Addr:       0x7fa3b88bbb30,
			Type:       "struct _opaque_pthread_t *",
			Kind:       "foo",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
