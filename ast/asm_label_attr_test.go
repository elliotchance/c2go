package ast

import (
	"testing"
)

func TestAsmLabelAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7ff26d8224e8 </usr/include/sys/cdefs.h:569:36> "_fopen"`: &AsmLabelAttr{
			Address:      "0x7ff26d8224e8",
			Position:     "/usr/include/sys/cdefs.h:569:36",
			Inherited:    false,
			FunctionName: "_fopen",
			Children:     []Node{},
		},
		`0x7fd55a169318 </usr/include/stdio.h:325:47> Inherited "_popen"`: &AsmLabelAttr{
			Address:      "0x7fd55a169318",
			Position:     "/usr/include/stdio.h:325:47",
			Inherited:    true,
			FunctionName: "_popen",
			Children:     []Node{},
		},
	}

	runNodeTests(t, nodes)
}
