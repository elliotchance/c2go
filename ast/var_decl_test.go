package ast

import (
	"testing"
)

func TestVarDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x7fd5e90e5a00 <col:14> col:17 'int'`: &VarDecl{
			Addr:         0x7fd5e90e5a00,
			Pos:          NewPositionFromString("col:14"),
			Position2:    "col:17",
			Name:         "",
			Type:         "int",
			Type2:        "",
			IsExtern:     false,
			IsUsed:       false,
			IsCInit:      false,
			IsReferenced: false,
			ChildNodes:   []Node{},
		},
		`0x7fd5e90e9078 <line:156:1, col:14> col:14 __stdinp 'FILE *' extern`: &VarDecl{
			Addr:         0x7fd5e90e9078,
			Pos:          NewPositionFromString("line:156:1, col:14"),
			Position2:    "col:14",
			Name:         "__stdinp",
			Type:         "FILE *",
			Type2:        "",
			IsExtern:     true,
			IsUsed:       false,
			IsCInit:      false,
			IsReferenced: false,
			ChildNodes:   []Node{},
		},
		`0x7fd5e90ed630 <col:40, col:47> col:47 __size 'size_t':'unsigned long'`: &VarDecl{
			Addr:         0x7fd5e90ed630,
			Pos:          NewPositionFromString("col:40, col:47"),
			Position2:    "col:47",
			Name:         "__size",
			Type:         "size_t",
			Type2:        "unsigned long",
			IsExtern:     false,
			IsUsed:       false,
			IsCInit:      false,
			IsReferenced: false,
			ChildNodes:   []Node{},
		},
		`0x7fee35907a78 <col:4, col:8> col:8 used c 'int'`: &VarDecl{
			Addr:         0x7fee35907a78,
			Pos:          NewPositionFromString("col:4, col:8"),
			Position2:    "col:8",
			Name:         "c",
			Type:         "int",
			Type2:        "",
			IsExtern:     false,
			IsUsed:       true,
			IsCInit:      false,
			IsReferenced: false,
			ChildNodes:   []Node{},
		},
		`0x7fb0fd90ba30 <col:3, /usr/include/sys/_types.h:52:33> tests/assert/assert.c:13:9 used b 'int *' cinit`: &VarDecl{
			Addr:         0x7fb0fd90ba30,
			Pos:          NewPositionFromString("col:3, /usr/include/sys/_types.h:52:33"),
			Position2:    "tests/assert/assert.c:13:9",
			Name:         "b",
			Type:         "int *",
			Type2:        "",
			IsExtern:     false,
			IsUsed:       true,
			IsCInit:      true,
			IsReferenced: false,
			ChildNodes:   []Node{},
		},
		`0x7fb20308bd40 <col:5, col:11> col:11 referenced a 'short'`: &VarDecl{
			Addr:         0x7fb20308bd40,
			Pos:          NewPositionFromString("col:5, col:11"),
			Position2:    "col:11",
			Name:         "a",
			Type:         "short",
			Type2:        "",
			IsExtern:     false,
			IsUsed:       false,
			IsCInit:      false,
			IsReferenced: true,
			ChildNodes:   []Node{},
		},
	}

	runNodeTests(t, nodes)
}
