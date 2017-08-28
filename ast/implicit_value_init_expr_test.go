package ast

import (
	"testing"
)

func TestImplicitValueInitExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7f8c3396fbd8 <<invalid sloc>> 'sqlite3StatValueType':'long long'`: &ImplicitValueInitExpr{
			Addr:       0x7f8c3396fbd8,
			Pos:        NewPositionFromString("<invalid sloc>"),
			Type1:      "sqlite3StatValueType",
			Type2:      "long long",
			ChildNodes: []Node{},
		},
		`0x7feecb0d6af0 <<invalid sloc>> 'char'`: &ImplicitValueInitExpr{
			Addr:       0x7feecb0d6af0,
			Pos:        NewPositionFromString("<invalid sloc>"),
			Type1:      "char",
			Type2:      "",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
