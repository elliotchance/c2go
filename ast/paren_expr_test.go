package ast

import (
	"testing"
)

func TestParenExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fb0bc8b2308 <col:10, col:25> 'unsigned char'`: &ParenExpr{
			Address:  "0x7fb0bc8b2308",
			Position: "col:10, col:25",
			Type:     "unsigned char",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
