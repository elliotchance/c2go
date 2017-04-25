package ast

import (
	"testing"
)

func TestAsmLabelAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7ff26d8224e8 </usr/include/sys/cdefs.h:569:36> "_fopen"`: &AsmLabelAttr{
			Address:      "0x7ff26d8224e8",
			Position:     "/usr/include/sys/cdefs.h:569:36",
			FunctionName: "_fopen",
			Children:     []Node{},
		},
	}

	runNodeTests(t, nodes)
}
