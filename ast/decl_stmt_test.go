package ast

import (
	"testing"
)

func TestDeclStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fb791846e80 <line:11:4, col:31>`: &DeclStmt{
			Address:  "0x7fb791846e80",
			Position: "line:11:4, col:31",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
