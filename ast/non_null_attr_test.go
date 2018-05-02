package ast

import (
	"testing"
)

func TestNonNullAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa1488273b0 <line:7:4, line:11:4> 1`: &NonNullAttr{
			Addr:       0x7fa1488273b0,
			Pos:        NewPositionFromString("line:7:4, line:11:4"),
			Inherited:  false,
			A:          1,
			B:          0,
			C:          0,
			D:          0,
			ChildNodes: []Node{},
		},
		`0x2cce280 </sys/cdefs.h:286:44, /bits/mathcalls.h:115:69> 1`: &NonNullAttr{
			Addr:       0x2cce280,
			Pos:        NewPositionFromString("/sys/cdefs.h:286:44, /bits/mathcalls.h:115:69"),
			Inherited:  false,
			A:          1,
			B:          0,
			C:          0,
			D:          0,
			ChildNodes: []Node{},
		},
		`0x201ede0 <line:145:79, col:93> 0`: &NonNullAttr{
			Addr:       0x201ede0,
			Pos:        NewPositionFromString("line:145:79, col:93"),
			Inherited:  false,
			A:          0,
			B:          0,
			C:          0,
			D:          0,
			ChildNodes: []Node{},
		},
		`0x1b89b20 <col:76, col:93> 2 3`: &NonNullAttr{
			Addr:       0x1b89b20,
			Pos:        NewPositionFromString("col:76, col:93"),
			Inherited:  false,
			A:          2,
			B:          3,
			C:          0,
			D:          0,
			ChildNodes: []Node{},
		},
		`0x55f0219e20d0 <line:717:22, col:42> 0 1 4`: &NonNullAttr{
			Addr:       0x55f0219e20d0,
			Pos:        NewPositionFromString("line:717:22, col:42"),
			Inherited:  false,
			A:          0,
			B:          1,
			C:          4,
			D:          0,
			ChildNodes: []Node{},
		},
		`0x248ea60 <line:155:26, col:49> 0 1 2 4`: &NonNullAttr{
			Addr:       0x248ea60,
			Pos:        NewPositionFromString("line:155:26, col:49"),
			Inherited:  false,
			A:          0,
			B:          1,
			C:          2,
			D:          4,
			ChildNodes: []Node{},
		},
		`0x39cf2b0 <col:53> Inherited 0 1`: &NonNullAttr{
			Addr:       0x39cf2b0,
			Pos:        NewPositionFromString("col:53"),
			Inherited:  true,
			A:          0,
			B:          1,
			C:          0,
			D:          0,
			ChildNodes: []Node{},
		},
		`0x2c3d600 <line:304:19>`: &NonNullAttr{
			Addr:       0x2c3d600,
			Pos:        NewPositionFromString("line:304:19"),
			Inherited:  false,
			A:          0,
			B:          0,
			C:          0,
			D:          0,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
