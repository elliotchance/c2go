package ast

import (
	"testing"
)

func TestTranslationUnitDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x7fe78a815ed0 <<invalid sloc>> <invalid sloc>`: &TranslationUnitDecl{
			Address:  "0x7fe78a815ed0",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
