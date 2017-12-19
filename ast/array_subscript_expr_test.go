package ast

import (
	"testing"
)

func TestArraySubscriptExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fe35b85d180 <col:63, col:69> 'char *' lvalue`: &ArraySubscriptExpr{
			Addr:       0x7fe35b85d180,
			Pos:        NewPositionFromString("col:63, col:69"),
			Type:       "char *",
			Type2:      "",
			IsLvalue:   true,
			ChildNodes: []Node{},
		},
		`0x2416660 <col:2, col:5> 'u32':'unsigned int' lvalue`: &ArraySubscriptExpr{
			Addr:       0x2416660,
			Pos:        NewPositionFromString("col:2, col:5"),
			Type:       "u32",
			Type2:      "unsigned int",
			IsLvalue:   true,
			ChildNodes: []Node{},
		},
		`0x3f147c0 <col:39, col:55> 'extCoord':'extCoord' lvalue`: &ArraySubscriptExpr{
			Addr:       0x3f147c0,
			Pos:        NewPositionFromString("col:39, col:55"),
			Type:       "extCoord",
			Type2:      "extCoord",
			IsLvalue:   true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
