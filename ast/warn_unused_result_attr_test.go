package ast

import (
	"testing"
)

func TestWarnUnusedResultAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa1d704d420 <col:60> warn_unused_result`: &WarnUnusedResultAttr{
			Addr:       0x7fa1d704d420,
			Pos:        NewPositionFromString("col:60"),
			Inherited:  false,
			ChildNodes: []Node{},
		},
		`0x1fac810 <line:481:52>`: &WarnUnusedResultAttr{
			Addr:       0x1fac810,
			Pos:        NewPositionFromString("line:481:52"),
			Inherited:  false,
			ChildNodes: []Node{},
		},
		`0x3794590 </home/kph/co/util-linux/libblkid/src/blkidP.h:374:19> Inherited`: &WarnUnusedResultAttr{
			Addr:       0x3794590,
			Pos:        NewPositionFromString("/home/kph/co/util-linux/libblkid/src/blkidP.h:374:19"),
			Inherited:  true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
