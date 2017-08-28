package ast

import (
	"testing"
)

func TestAlwaysInlineAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fce780f5018 </usr/include/sys/cdefs.h:313:68> always_inline`: &AlwaysInlineAttr{
			Addr:       0x7fce780f5018,
			Pos:        NewPositionFromString("/usr/include/sys/cdefs.h:313:68"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
