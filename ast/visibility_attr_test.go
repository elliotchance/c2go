package ast

import (
	"testing"
)

func TestVisibilityAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x55c49d8dd1d8 <col:16, col:36> Default`: &VisibilityAttr{
			Addr:        0x55c49d8dd1d8,
			Pos:         NewPositionFromString("col:16, col:36"),
			ChildNodes:  []Node{},
			IsInherited: false,
			IsDefault:   true,
			IsHidden:    false,
		},
		`0x7f8e7b00bb80 </cmark/src/cmark.h:497:16, col:36> Inherited Default`: &VisibilityAttr{
			Addr:        0x7f8e7b00bb80,
			Pos:         NewPositionFromString("/cmark/src/cmark.h:497:16, col:36"),
			ChildNodes:  []Node{},
			IsInherited: true,
			IsDefault:   true,
			IsHidden:    false,
		},
		`0x55ab30581650 <line:24:16, col:35> Hidden`: &VisibilityAttr{
			Addr:        0x55ab30581650,
			Pos:         NewPositionFromString("line:24:16, col:35"),
			ChildNodes:  []Node{},
			IsInherited: false,
			IsDefault:   false,
			IsHidden:    true,
		},
	}

	runNodeTests(t, nodes)
}
