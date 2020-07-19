package ast

import (
	"testing"
)

func TestAsmLabelAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7ff26d8224e8 </usr/include/sys/cdefs.h:569:36> "_fopen"`: &AsmLabelAttr{
			Addr:           0x7ff26d8224e8,
			Pos:            NewPositionFromString("/usr/include/sys/cdefs.h:569:36"),
			Inherited:      false,
			FunctionName:   "_fopen",
			ChildNodes:     []Node{},
			IsLiteralLabel: false,
		},
		`0x7fd55a169318 </usr/include/stdio.h:325:47> Inherited "_popen"`: &AsmLabelAttr{
			Addr:           0x7fd55a169318,
			Pos:            NewPositionFromString("/usr/include/stdio.h:325:47"),
			Inherited:      true,
			FunctionName:   "_popen",
			ChildNodes:     []Node{},
			IsLiteralLabel: false,
		},
		`0x559fea32f5f0 <line:407:94> "__isoc99_fscanf" IsLiteralLabel`: &AsmLabelAttr{
			Addr:           0x559fea32f5f0,
			Pos:            NewPositionFromString("line:407:94"),
			Inherited:      false,
			FunctionName:   "__isoc99_fscanf",
			ChildNodes:     []Node{},
			IsLiteralLabel: true,
		},
	}

	runNodeTests(t, nodes)
}
