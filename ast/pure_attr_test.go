package ast

import (
	"testing"
)

func TestPureAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fe9eb899198 <col:1> Implicit`: &PureAttr{
			Address:   "0x7fe9eb899198",
			Position:  "col:1",
			Implicit:  true,
			Inherited: false,
			Children:  []interface{}{},
		},
		`0x7fe8d60992a0 <col:1> Inherited Implicit`: &PureAttr{
			Address:   "0x7fe8d60992a0",
			Position:  "col:1",
			Implicit:  true,
			Inherited: true,
			Children:  []interface{}{},
		},
	}

	runNodeTests(t, nodes)
}
