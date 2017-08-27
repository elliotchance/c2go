package ast

import (
	"testing"
)

func TestWarnUnusedResultAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa1d704d420 <col:60> warn_unused_result`: &WarnUnusedResultAttr{
			Addr:     0x7fa1d704d420,
			Position: "col:60",
			Children: []Node{},
		},
		`0x1fac810 <line:481:52>`: &WarnUnusedResultAttr{
			Addr:     0x1fac810,
			Position: "line:481:52",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
