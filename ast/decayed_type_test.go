package ast

import (
	"testing"
)

func TestDecayedType(t *testing.T) {
	nodes := map[string]Node{
		`0x7f1234567890 'struct __va_list_tag *' sugar`: &DecayedType{
			Addr:       0x7f1234567890,
			Type:       "struct __va_list_tag *",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
