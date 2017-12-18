package ast

import (
	"testing"
)

func TestSentinelAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x346df70 <line:3571:19, col:33> 0 0`: &SentinelAttr{
			Addr:       0x346df70,
			Pos:        NewPositionFromString("line:3571:19, col:33"),
			A:          0,
			B:          0,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
