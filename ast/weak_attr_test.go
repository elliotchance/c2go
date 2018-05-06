package ast

import (
	"testing"
)

func TestWeakAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x56069ece5110 <line:736:22>`: &WeakAttr{
			Addr:       0x56069ece5110,
			Pos:        NewPositionFromString("line:736:22"),
			Inherited:  false,
			ChildNodes: []Node{},
		},
		`0x20c6ad0 </glibc-2.27/support/temp_file-internal.h:27:62> Inherited`: &WeakAttr{
			Addr:       0x20c6ad0,
			Pos:        NewPositionFromString("/glibc-2.27/support/temp_file-internal.h:27:62"),
			Inherited:  true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
