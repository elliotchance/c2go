package ast

import (
	"testing"
)

func TestRecord(t *testing.T) {
	nodes := map[string]Node{
		`0x7fd3ab857950 '__sFILE'`: &Record{
			Address:  "0x7fd3ab857950",
			Type:     "__sFILE",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
