package ast

import (
	"testing"
)

func TestParmVarDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x7f973380f000 <col:14> col:17 'int'`: &ParmVarDecl{
			Addr:       0x7f973380f000,
			Pos:        NewPositionFromString("col:14"),
			Position2:  "col:17",
			Type:       "int",
			Name:       "",
			Type2:      "",
			IsUsed:     false,
			ChildNodes: []Node{},
		},
		`0x7f973380f070 <col:19, col:30> col:31 'const char *'`: &ParmVarDecl{
			Addr:       0x7f973380f070,
			Pos:        NewPositionFromString("col:19, col:30"),
			Position2:  "col:31",
			Type:       "const char *",
			Name:       "",
			Type2:      "",
			IsUsed:     false,
			ChildNodes: []Node{},
		},
		`0x7f9733816e50 <col:13, col:37> col:37 __filename 'const char *__restrict'`: &ParmVarDecl{
			Addr:       0x7f9733816e50,
			Pos:        NewPositionFromString("col:13, col:37"),
			Position2:  "col:37",
			Type:       "const char *__restrict",
			Name:       "__filename",
			Type2:      "",
			IsUsed:     false,
			ChildNodes: []Node{},
		},
		`0x7f9733817418 <<invalid sloc>> <invalid sloc> 'FILE *'`: &ParmVarDecl{
			Addr:       0x7f9733817418,
			Pos:        NewPositionFromString("<invalid sloc>"),
			Position2:  "<invalid sloc>",
			Type:       "FILE *",
			Name:       "",
			Type2:      "",
			IsUsed:     false,
			ChildNodes: []Node{},
		},
		`0x7f9733817c30 <col:40, col:47> col:47 __size 'size_t':'unsigned long'`: &ParmVarDecl{
			Addr:       0x7f9733817c30,
			Pos:        NewPositionFromString("col:40, col:47"),
			Position2:  "col:47",
			Type:       "size_t",
			Name:       "__size",
			Type2:      "unsigned long",
			IsUsed:     false,
			ChildNodes: []Node{},
		},
		`0x7f973382fa10 <line:476:18, col:25> col:34 'int (* _Nullable)(void *, char *, int)':'int (*)(void *, char *, int)'`: &ParmVarDecl{
			Addr:       0x7f973382fa10,
			Pos:        NewPositionFromString("line:476:18, col:25"),
			Position2:  "col:34",
			Type:       "int (* _Nullable)(void *, char *, int)",
			Name:       "",
			Type2:      "int (*)(void *, char *, int)",
			IsUsed:     false,
			ChildNodes: []Node{},
		},
		`0x7f97338355b8 <col:10, col:14> col:14 used argc 'int'`: &ParmVarDecl{
			Addr:       0x7f97338355b8,
			Pos:        NewPositionFromString("col:10, col:14"),
			Position2:  "col:14",
			Type:       "int",
			Name:       "argc",
			Type2:      "",
			IsUsed:     true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
