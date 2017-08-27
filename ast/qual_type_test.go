package ast

import (
	"testing"
)

func TestQualType(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa3b88bbb31 'struct _opaque_pthread_t *' foo`: &QualType{
			Addr:       0x7fa3b88bbb31,
			Type:       "struct _opaque_pthread_t *",
			Kind:       "foo",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
