package ast

import (
	"testing"
)

func TestFieldDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x7fef510c4848 <line:141:2, col:6> col:6 _ur 'int'`: &FieldDecl{
			Addr:       0x7fef510c4848,
			Pos:        NewPositionFromString("line:141:2, col:6"),
			Position2:  "col:6",
			Name:       "_ur",
			Type:       "int",
			Implicit:   false,
			Referenced: false,
			ChildNodes: []Node{},
		},
		`0x7fef510c46f8 <line:139:2, col:16> col:16 _ub 'struct __sbuf':'struct __sbuf'`: &FieldDecl{
			Addr:       0x7fef510c46f8,
			Pos:        NewPositionFromString("line:139:2, col:16"),
			Position2:  "col:16",
			Name:       "_ub",
			Type:       "struct __sbuf",
			Implicit:   false,
			Referenced: false,
			ChildNodes: []Node{},
		},
		`0x7fef510c3fe0 <line:134:2, col:19> col:19 _read 'int (* _Nullable)(void *, char *, int)':'int (*)(void *, char *, int)'`: &FieldDecl{
			Addr:       0x7fef510c3fe0,
			Pos:        NewPositionFromString("line:134:2, col:19"),
			Position2:  "col:19",
			Name:       "_read",
			Type:       "int (* _Nullable)(void *, char *, int)",
			Implicit:   false,
			Referenced: false,
			ChildNodes: []Node{},
		},
		`0x7fef51073a60 <line:105:2, col:40> col:40 __cleanup_stack 'struct __darwin_pthread_handler_rec *'`: &FieldDecl{
			Addr:       0x7fef51073a60,
			Pos:        NewPositionFromString("line:105:2, col:40"),
			Position2:  "col:40",
			Name:       "__cleanup_stack",
			Type:       "struct __darwin_pthread_handler_rec *",
			Implicit:   false,
			Referenced: false,
			ChildNodes: []Node{},
		},
		`0x7fef510738e8 <line:100:2, col:43> col:7 __opaque 'char [16]'`: &FieldDecl{
			Addr:       0x7fef510738e8,
			Pos:        NewPositionFromString("line:100:2, col:43"),
			Position2:  "col:7",
			Name:       "__opaque",
			Type:       "char [16]",
			Implicit:   false,
			Referenced: false,
			ChildNodes: []Node{},
		},
		`0x7fe9f5072268 <line:129:2, col:6> col:6 referenced _lbfsize 'int'`: &FieldDecl{
			Addr:       0x7fe9f5072268,
			Pos:        NewPositionFromString("line:129:2, col:6"),
			Position2:  "col:6",
			Name:       "_lbfsize",
			Type:       "int",
			Implicit:   false,
			Referenced: true,
			ChildNodes: []Node{},
		},
		`0x7f9bc9083d00 <line:91:5, line:97:8> line:91:5 'unsigned short'`: &FieldDecl{
			Addr:       0x7f9bc9083d00,
			Pos:        NewPositionFromString("line:91:5, line:97:8"),
			Position2:  "line:91:5",
			Name:       "",
			Type:       "unsigned short",
			Implicit:   false,
			Referenced: false,
			ChildNodes: []Node{},
		},
		`0x30363a0 <col:18, col:29> __val 'int [2]'`: &FieldDecl{
			Addr:       0x30363a0,
			Pos:        NewPositionFromString("col:18, col:29"),
			Position2:  "",
			Name:       "__val",
			Type:       "int [2]",
			Implicit:   false,
			Referenced: false,
			ChildNodes: []Node{},
		},
		`0x17aeac0 <line:3:9> col:9 implicit referenced 'struct vec3d_t::(anonymous at main.c:3:9)'`: &FieldDecl{
			Addr:       0x17aeac0,
			Pos:        NewPositionFromString("line:3:9"),
			Position2:  "col:9",
			Name:       "",
			Type:       "struct vec3d_t::(anonymous at main.c:3:9)",
			Implicit:   true,
			Referenced: true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
