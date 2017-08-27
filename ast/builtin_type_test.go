package ast

import (
	"testing"
)

func TestBuiltinType(t *testing.T) {
	nodes := map[string]Node{
		`0x7f8a43023f40 '__int128'`: &BuiltinType{
			Addr:       0x7f8a43023f40,
			Type:       "__int128",
			ChildNodes: []Node{},
		},
		`0x7f8a43023ea0 'unsigned long long'`: &BuiltinType{
			Addr:       0x7f8a43023ea0,
			Type:       "unsigned long long",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
