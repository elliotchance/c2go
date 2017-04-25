package ast

import (
	"testing"
)

func TestFunctionDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x7fb5a90e60d0 <line:231:1, col:22> col:7 clearerr 'void (FILE *)'`: &FunctionDecl{
			Address:    "0x7fb5a90e60d0",
			Position:   "line:231:1, col:22",
			Prev:       "",
			Position2:  "col:7",
			Name:       "clearerr",
			Type:       "void (FILE *)",
			IsExtern:   false,
			IsImplicit: false,
			IsUsed:     false,
			Children:   []Node{},
		},
		`0x7fb5a90e2a50 </usr/include/sys/stdio.h:39:1, /usr/include/AvailabilityInternal.h:21697:126> /usr/include/sys/stdio.h:39:5 renameat 'int (int, const char *, int, const char *)'`: &FunctionDecl{
			Address:    "0x7fb5a90e2a50",
			Position:   "/usr/include/sys/stdio.h:39:1, /usr/include/AvailabilityInternal.h:21697:126",
			Prev:       "",
			Position2:  "/usr/include/sys/stdio.h:39:5",
			Name:       "renameat",
			Type:       "int (int, const char *, int, const char *)",
			IsExtern:   false,
			IsImplicit: false,
			IsUsed:     false,
			Children:   []Node{},
		},
		`0x7fb5a90e9b70 </usr/include/stdio.h:244:6> col:6 implicit fprintf 'int (FILE *, const char *, ...)' extern`: &FunctionDecl{
			Address:    "0x7fb5a90e9b70",
			Position:   "/usr/include/stdio.h:244:6",
			Prev:       "",
			Position2:  "col:6",
			Name:       "fprintf",
			Type:       "int (FILE *, const char *, ...)",
			IsExtern:   true,
			IsImplicit: true,
			IsUsed:     false,
			Children:   []Node{},
		},
		`0x7fb5a90e9d40 prev 0x7fb5a90e9b70 <col:1, /usr/include/sys/cdefs.h:351:63> /usr/include/stdio.h:244:6 fprintf 'int (FILE *, const char *, ...)'`: &FunctionDecl{
			Address:    "0x7fb5a90e9d40",
			Position:   "col:1, /usr/include/sys/cdefs.h:351:63",
			Prev:       "0x7fb5a90e9b70",
			Position2:  "/usr/include/stdio.h:244:6",
			Name:       "fprintf",
			Type:       "int (FILE *, const char *, ...)",
			IsExtern:   false,
			IsImplicit: false,
			IsUsed:     false,
			Children:   []Node{},
		},
		`0x7fb5a90ec210 <line:259:6> col:6 implicit used printf 'int (const char *, ...)' extern`: &FunctionDecl{
			Address:    "0x7fb5a90ec210",
			Position:   "line:259:6",
			Prev:       "",
			Position2:  "col:6",
			Name:       "printf",
			Type:       "int (const char *, ...)",
			IsExtern:   true,
			IsImplicit: true,
			IsUsed:     true,
			Children:   []Node{},
		},
		`0x2ae30d8 </usr/include/math.h:65:3, /usr/include/x86_64-linux-gnu/sys/cdefs.h:57:54> <scratch space>:17:1 __acos 'double (double)' extern`: &FunctionDecl{
			Address:    "0x2ae30d8",
			Position:   "/usr/include/math.h:65:3, /usr/include/x86_64-linux-gnu/sys/cdefs.h:57:54",
			Prev:       "",
			Position2:  "<scratch space>:17:1",
			Name:       "__acos",
			Type:       "double (double)",
			IsExtern:   true,
			IsImplicit: false,
			IsUsed:     false,
			Children:   []Node{},
		},
	}

	runNodeTests(t, nodes)
}
