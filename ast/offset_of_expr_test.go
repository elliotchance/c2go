package ast

import (
	"testing"
)

func TestOffsetOfExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa855aab838 <col:63, col:95> 'unsigned long'`: &OffsetOfExpr{
			Addr:       0x7fa855aab838,
			Pos:        NewPositionFromString("col:63, col:95"),
			Type:       "unsigned long",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
