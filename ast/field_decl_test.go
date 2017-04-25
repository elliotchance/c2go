package ast

import (
	"testing"
)

func TestFieldDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x7fef510c4848 <line:141:2, col:6> col:6 _ur 'int'`: &FieldDecl{
			Address:    "0x7fef510c4848",
			Position:   "line:141:2, col:6",
			Position2:  "col:6",
			Name:       "_ur",
			Type:       "int",
			Referenced: false,
			Children:   []Node{},
		},
		`0x7fef510c46f8 <line:139:2, col:16> col:16 _ub 'struct __sbuf':'struct __sbuf'`: &FieldDecl{
			Address:    "0x7fef510c46f8",
			Position:   "line:139:2, col:16",
			Position2:  "col:16",
			Name:       "_ub",
			Type:       "struct __sbuf",
			Referenced: false,
			Children:   []Node{},
		},
		`0x7fef510c3fe0 <line:134:2, col:19> col:19 _read 'int (* _Nullable)(void *, char *, int)':'int (*)(void *, char *, int)'`: &FieldDecl{
			Address:    "0x7fef510c3fe0",
			Position:   "line:134:2, col:19",
			Position2:  "col:19",
			Name:       "_read",
			Type:       "int (* _Nullable)(void *, char *, int)",
			Referenced: false,
			Children:   []Node{},
		},
		`0x7fef51073a60 <line:105:2, col:40> col:40 __cleanup_stack 'struct __darwin_pthread_handler_rec *'`: &FieldDecl{
			Address:    "0x7fef51073a60",
			Position:   "line:105:2, col:40",
			Position2:  "col:40",
			Name:       "__cleanup_stack",
			Type:       "struct __darwin_pthread_handler_rec *",
			Referenced: false,
			Children:   []Node{},
		},
		`0x7fef510738e8 <line:100:2, col:43> col:7 __opaque 'char [16]'`: &FieldDecl{
			Address:    "0x7fef510738e8",
			Position:   "line:100:2, col:43",
			Position2:  "col:7",
			Name:       "__opaque",
			Type:       "char [16]",
			Referenced: false,
			Children:   []Node{},
		},
		`0x7fe9f5072268 <line:129:2, col:6> col:6 referenced _lbfsize 'int'`: &FieldDecl{
			Address:    "0x7fe9f5072268",
			Position:   "line:129:2, col:6",
			Position2:  "col:6",
			Name:       "_lbfsize",
			Type:       "int",
			Referenced: true,
			Children:   []Node{},
		},
		`0x7f9bc9083d00 <line:91:5, line:97:8> line:91:5 'unsigned short'`: &FieldDecl{
			Address:    "0x7f9bc9083d00",
			Position:   "line:91:5, line:97:8",
			Position2:  "line:91:5",
			Name:       "",
			Type:       "unsigned short",
			Referenced: false,
			Children:   []Node{},
		},
		`0x30363a0 <col:18, col:29> __val 'int [2]'`: &FieldDecl{
			Address:    "0x30363a0",
			Position:   "col:18, col:29",
			Position2:  "",
			Name:       "__val",
			Type:       "int [2]",
			Referenced: false,
			Children:   []Node{},
		},
	}

	runNodeTests(t, nodes)
}
