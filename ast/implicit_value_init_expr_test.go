package ast

import (
	"testing"
)

func TestImplicitValueInitExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7f8c3396fbd8 <<invalid sloc>> 'sqlite3StatValueType':'long long'`: &ImplicitValueInitExpr{
			Address:  "0x7f8c3396fbd8",
			Position: "<invalid sloc>",
			Type1:    "sqlite3StatValueType",
			Type2:    "long long",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
