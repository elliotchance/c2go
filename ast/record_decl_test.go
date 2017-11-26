package ast

import (
	"testing"
)

func TestRecordDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x7f913c0dbb50 <line:76:9, line:79:1> line:76:9 union definition`: &RecordDecl{
			Addr:       0x7f913c0dbb50,
			Pos:        NewPositionFromString("line:76:9, line:79:1"),
			Prev:       "",
			Position2:  "line:76:9",
			Kind:       "union",
			Name:       "",
			Definition: true,
			ChildNodes: []Node{},
		},
		`0x7f85360285c8 </usr/include/sys/_pthread/_pthread_types.h:57:1, line:61:1> line:57:8 struct __darwin_pthread_handler_rec definition`: &RecordDecl{
			Addr:       0x7f85360285c8,
			Pos:        NewPositionFromString("/usr/include/sys/_pthread/_pthread_types.h:57:1, line:61:1"),
			Prev:       "",
			Position2:  "line:57:8",
			Kind:       "struct",
			Name:       "__darwin_pthread_handler_rec",
			Definition: true,
			ChildNodes: []Node{},
		},
		`0x7f85370248a0 <line:94:1, col:8> col:8 struct __sFILEX`: &RecordDecl{
			Addr:       0x7f85370248a0,
			Pos:        NewPositionFromString("line:94:1, col:8"),
			Prev:       "",
			Position2:  "col:8",
			Kind:       "struct",
			Name:       "__sFILEX",
			Definition: false,
			ChildNodes: []Node{},
		},
		`0x5564ed488a10 parent 0x5564ed3ffe00 <line:7232:3, line:7237:3> line:7232:10 struct sqlite3_index_constraint definition`: &RecordDecl{
			Addr:       0x5564ed488a10,
			Pos:        NewPositionFromString("line:7232:3, line:7237:3"),
			Prev:       "",
			Position2:  "line:7232:10",
			Kind:       "struct",
			Name:       "sqlite3_index_constraint",
			Definition: true,
			ChildNodes: []Node{},
		},
		`0x56454e55e4b8 prev 0x56454e55e360 <line:86428:1, line:86437:1> line:86428:8 struct Incrblob definition`: &RecordDecl{
			Addr:       0x56454e55e4b8,
			Pos:        NewPositionFromString("line:86428:1, line:86437:1"),
			Prev:       "0x56454e55e360",
			Position2:  "line:86428:8",
			Kind:       "struct",
			Name:       "Incrblob",
			Definition: true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
