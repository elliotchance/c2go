package ast

import (
	"testing"
)

func TestNonNullAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa1488273b0 <line:7:4, line:11:4> 1`: &NonNullAttr{
			Addr:       0x7fa1488273b0,
			Pos:        "line:7:4, line:11:4",
			A:          1,
			B:          0,
			ChildNodes: []Node{},
		},
		`0x2cce280 </sys/cdefs.h:286:44, /bits/mathcalls.h:115:69> 1`: &NonNullAttr{
			Addr:       0x2cce280,
			Pos:        "/sys/cdefs.h:286:44, /bits/mathcalls.h:115:69",
			A:          1,
			B:          0,
			ChildNodes: []Node{},
		},
		`0x201ede0 <line:145:79, col:93> 0`: &NonNullAttr{
			Addr:       0x201ede0,
			Pos:        "line:145:79, col:93",
			A:          0,
			B:          0,
			ChildNodes: []Node{},
		},
		`0x1b89b20 <col:76, col:93> 2 3`: &NonNullAttr{
			Addr:       0x1b89b20,
			Pos:        "col:76, col:93",
			A:          2,
			B:          3,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
