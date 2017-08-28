package ast

import (
	"testing"
)

func TestCharacterLiteral(t *testing.T) {
	nodes := map[string]Node{
		`0x7f980b858308 <col:62> 'int' 10`: &CharacterLiteral{
			Addr:       0x7f980b858308,
			Pos:        NewPositionFromString("col:62"),
			Type:       "int",
			Value:      10,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
