package ast

import (
	"testing"
)

func TestFloatingLiteral(t *testing.T) {
	nodes := map[string]Node{
		`0x7febe106f5e8 <col:24> 'double' 1.230000e+00`: &FloatingLiteral{
			Address:  "0x7febe106f5e8",
			Position: "col:24",
			Type:     "double",
			Value:    1.23,
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
