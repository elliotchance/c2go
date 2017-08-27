package ast

import (
	"testing"
)

func TestParenType(t *testing.T) {
	nodes := map[string]Node{
		`0x7faf820a4c60 'void (int)' sugar`: &ParenType{
			Addr:     0x7faf820a4c60,
			Type:     "void (int)",
			Sugar:    true,
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
