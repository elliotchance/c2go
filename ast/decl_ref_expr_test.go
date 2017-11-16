package ast

import (
	"testing"
)

func TestDeclRefExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fc972064460 <col:8> 'FILE *' lvalue ParmVar 0x7fc9720642d0 '_p' 'FILE *'`: &DeclRefExpr{
			Addr:       0x7fc972064460,
			Pos:        NewPositionFromString("col:8"),
			Type:       "FILE *",
			Type1:      "",
			Lvalue:     true,
			For:        "ParmVar",
			Address2:   "0x7fc9720642d0",
			Name:       "_p",
			Type2:      "FILE *",
			Type3:      "",
			ChildNodes: []Node{},
		},
		`0x7fc97206a958 <col:11> 'int (int, FILE *)' Function 0x7fc972064198 '__swbuf' 'int (int, FILE *)'`: &DeclRefExpr{
			Addr:       0x7fc97206a958,
			Pos:        NewPositionFromString("col:11"),
			Type:       "int (int, FILE *)",
			Type1:      "",
			Lvalue:     false,
			For:        "Function",
			Address2:   "0x7fc972064198",
			Name:       "__swbuf",
			Type2:      "int (int, FILE *)",
			Type3:      "",
			ChildNodes: []Node{},
		},
		`0x7fa36680f170 <col:19> 'struct programming':'struct programming' lvalue Var 0x7fa36680dc20 'variable' 'struct programming':'struct programming'`: &DeclRefExpr{
			Addr:       0x7fa36680f170,
			Pos:        NewPositionFromString("col:19"),
			Type:       "struct programming",
			Type1:      "struct programming",
			Lvalue:     true,
			For:        "Var",
			Address2:   "0x7fa36680dc20",
			Name:       "variable",
			Type2:      "struct programming",
			Type3:      "struct programming",
			ChildNodes: []Node{},
		},
		`0x35cb438 <col:13> 'int' EnumConstant 0x35ca300 'Jan' 'int'`: &DeclRefExpr{
			Addr:       0x35cb438,
			Pos:        NewPositionFromString("col:13"),
			Type:       "int",
			Type1:      "",
			Lvalue:     false,
			For:        "EnumConstant",
			Address2:   "0x35ca300",
			Name:       "Jan",
			Type2:      "int",
			Type3:      "",
			ChildNodes: []Node{},
		},
		`0x1ff8770 <col:33> 'T_ENUM':'T_ENUM' lvalue Var 0x1ff8600 'cc' 'T_ENUM':'T_ENUM'`: &DeclRefExpr{
			Addr:       0x1ff8770,
			Pos:        NewPositionFromString("col:33"),
			Type:       "T_ENUM",
			Type1:      "T_ENUM",
			Lvalue:     true,
			For:        "Var",
			Address2:   "0x1ff8600",
			Name:       "cc",
			Type2:      "T_ENUM",
			Type3:      "T_ENUM",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
