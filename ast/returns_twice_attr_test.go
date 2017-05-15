package ast

import (
	"testing"
)

func TestReturnsTwiceAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7ff8e9091640 <col:7> Implicit`: &ReturnsTwiceAttr{
			Address:  "0x7ff8e9091640",
			Position: "col:7",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
