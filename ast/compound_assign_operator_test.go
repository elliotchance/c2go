package ast

import (
	"testing"
)

func TestCompoundAssignOperator(t *testing.T) {
	nodes := map[string]Node{
		`0x2dc5758 <line:5:2, col:7> 'int' '+=' ComputeLHSTy='int' ComputeResultTy='int'`: &CompoundAssignOperator{
			Addr:                  0x2dc5758,
			Pos:                   NewPositionFromString("line:5:2, col:7"),
			Type:                  "int",
			Opcode:                "+=",
			ComputationLHSType:    "int",
			ComputationResultType: "int",
			ChildNodes:            []Node{},
		},
	}

	runNodeTests(t, nodes)
}
