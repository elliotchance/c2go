package ast

import (
	"testing"
)

func TestTypedefDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x7fdef0862430 <line:120:1, col:16> col:16`: &TypedefDecl{
			Addr:         0x7fdef0862430,
			Pos:          NewPositionFromString("line:120:1, col:16"),
			Position2:    "col:16",
			Name:         "",
			Type:         "",
			Type2:        "",
			IsImplicit:   false,
			IsReferenced: false,
			ChildNodes:   []Node{},
		},
		`0x7ffb9f824278 <<invalid sloc>> <invalid sloc> implicit __uint128_t 'unsigned __int128'`: &TypedefDecl{
			Addr:         0x7ffb9f824278,
			Pos:          NewPositionFromString("<invalid sloc>"),
			Position2:    "<invalid sloc>",
			Name:         "__uint128_t",
			Type:         "unsigned __int128",
			Type2:        "",
			IsImplicit:   true,
			IsReferenced: false,
			ChildNodes:   []Node{},
		},
		`0x7ffb9f824898 <<invalid sloc>> <invalid sloc> implicit referenced __builtin_va_list 'struct __va_list_tag [1]'`: &TypedefDecl{
			Addr:         0x7ffb9f824898,
			Pos:          NewPositionFromString("<invalid sloc>"),
			Position2:    "<invalid sloc>",
			Name:         "__builtin_va_list",
			Type:         "struct __va_list_tag [1]",
			Type2:        "",
			IsImplicit:   true,
			IsReferenced: true,
			ChildNodes:   []Node{},
		},
		`0x7ffb9f8248f8 </usr/include/i386/_types.h:37:1, col:24> col:24 __int8_t 'signed char'`: &TypedefDecl{
			Addr:         0x7ffb9f8248f8,
			Pos:          NewPositionFromString("/usr/include/i386/_types.h:37:1, col:24"),
			Position2:    "col:24",
			Name:         "__int8_t",
			Type:         "signed char",
			Type2:        "",
			IsImplicit:   false,
			IsReferenced: false,
			ChildNodes:   []Node{},
		},
		`0x7ffb9f8dbf50 <line:98:1, col:27> col:27 referenced __darwin_va_list '__builtin_va_list':'struct __va_list_tag [1]'`: &TypedefDecl{
			Addr:         0x7ffb9f8dbf50,
			Pos:          NewPositionFromString("line:98:1, col:27"),
			Position2:    "col:27",
			Name:         "__darwin_va_list",
			Type:         "__builtin_va_list",
			Type2:        "struct __va_list_tag [1]",
			IsImplicit:   false,
			IsReferenced: true,
			ChildNodes:   []Node{},
		},
		`0x34461f0 <line:338:1, col:77> __io_read_fn '__ssize_t (void *, char *, size_t)'`: &TypedefDecl{
			Addr:         0x34461f0,
			Pos:          NewPositionFromString("line:338:1, col:77"),
			Position2:    "",
			Name:         "__io_read_fn",
			Type:         "__ssize_t (void *, char *, size_t)",
			Type2:        "",
			IsImplicit:   false,
			IsReferenced: false,
			ChildNodes:   []Node{},
		},
		`0x55b9da8784b0 <line:341:1, line:342:16> line:341:19 __io_write_fn '__ssize_t (void *, const char *, size_t)'`: &TypedefDecl{
			Addr:         0x55b9da8784b0,
			Pos:          NewPositionFromString("line:341:1, line:342:16"),
			Position2:    "line:341:19",
			Name:         "__io_write_fn",
			Type:         "__ssize_t (void *, const char *, size_t)",
			Type2:        "",
			IsImplicit:   false,
			IsReferenced: false,
			ChildNodes:   []Node{},
		},
		`0x3f0b9b0 <line:12:1, line:15:3> col:3 referenced extCoord 'struct extCoord':'extCoord'`: &TypedefDecl{
			Addr:         0x3f0b9b0,
			Pos:          NewPositionFromString("line:12:1, line:15:3"),
			Position2:    "col:3",
			Name:         "extCoord",
			Type:         "struct extCoord",
			Type2:        "extCoord",
			IsImplicit:   false,
			IsReferenced: true,
			ChildNodes:   []Node{},
		},
	}

	runNodeTests(t, nodes)
}
