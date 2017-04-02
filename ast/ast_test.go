package ast_test

import (
	"testing"
	"github.com/elliotchance/c2go/ast"
	"reflect"
)

var nodes = map[string]interface{}{
	// AlwaysInlineAttr
	`0x7fce780f5018 </usr/include/sys/cdefs.h:313:68> always_inline`:
	ast.AlwaysInlineAttr{
		Address: "0x7fce780f5018",
		Position: "/usr/include/sys/cdefs.h:313:68",
	},

	// ArraySubscriptExpr
	`0x7fe35b85d180 <col:63, col:69> 'char *' lvalue`:
	ast.ArraySubscriptExpr{
		Address: "0x7fe35b85d180",
		Position: "col:63, col:69",
		Type: "char *",
		Kind: "lvalue",
	},

	// AsmLabelAttr
	`0x7ff26d8224e8 </usr/include/sys/cdefs.h:569:36> "_fopen"`:
	ast.AsmLabelAttr{
		Address: "0x7ff26d8224e8",
		Position: "/usr/include/sys/cdefs.h:569:36",
		FunctionName: "_fopen",
	},

	// AvailabilityAttr
	`0x7fc5ff8e5d18 </usr/include/AvailabilityInternal.h:21697:88, col:124> macos 10.10 0 0 "" ""`:
	ast.AvailabilityAttr{
		Address: "0x7fc5ff8e5d18",
		Position: "/usr/include/AvailabilityInternal.h:21697:88, col:124",
		OS: "macos",
		Version: "10.10",
		Unknown1: 0,
		Unknown2: 0,
		Unavailable: false,
		Message1: "",
		Message2: "",
	},
	`0x7fc5ff8e60d0 </usr/include/Availability.h:215:81, col:115> watchos 3.0 0 0 "" ""`:
	ast.AvailabilityAttr{
		Address: "0x7fc5ff8e60d0",
		Position: "/usr/include/Availability.h:215:81, col:115",
		OS: "watchos",
		Version: "3.0",
		Unknown1: 0,
		Unknown2: 0,
		Unavailable: false,
		Message1: "",
		Message2: "",
	},
	`0x7fc5ff8e6170 <col:81, col:115> tvos 10.0 0 0 "" ""`:
	ast.AvailabilityAttr{
		Address: "0x7fc5ff8e6170",
		Position: "col:81, col:115",
		OS: "tvos",
		Version: "10.0",
		Unknown1: 0,
		Unknown2: 0,
		Unavailable: false,
		Message1: "",
		Message2: "",
	},
	`0x7fc5ff8e61d8 <col:81, col:115> ios 10.0 0 0 "" ""`:
	ast.AvailabilityAttr{
		Address: "0x7fc5ff8e61d8",
		Position: "col:81, col:115",
		OS: "ios",
		Version: "10.0",
		Unknown1: 0,
		Unknown2: 0,
		Unavailable: false,
		Message1: "",
		Message2: "",
	},
	`0x7fc5ff8f0e18 </usr/include/sys/cdefs.h:275:50, col:99> swift 0 0 0 Unavailable "Use snprintf instead." ""`:
	ast.AvailabilityAttr{
		Address: "0x7fc5ff8f0e18",
		Position: "/usr/include/sys/cdefs.h:275:50, col:99",
		OS: "swift",
		Version: "0",
		Unknown1: 0,
		Unknown2: 0,
		Unavailable: true,
		Message1: "Use snprintf instead.",
		Message2: "",
	},
	`0x7fc5ff8f1988 <line:275:50, col:99> swift 0 0 0 Unavailable "Use mkstemp(3) instead." ""`:
	ast.AvailabilityAttr{
		Address: "0x7fc5ff8f1988",
		Position: "line:275:50, col:99",
		OS: "swift",
		Version: "0",
		Unknown1: 0,
		Unknown2: 0,
		Unavailable: true,
		Message1: "Use mkstemp(3) instead.",
		Message2: "",
	},

	// BinaryOperator
	`0x7fca2d8070e0 <col:11, col:23> 'unsigned char' '='`:
	ast.BinaryOperator{
		Address: "0x7fca2d8070e0",
		Position: "col:11, col:23",
		Type: "unsigned char",
		Operator: "=",
	},

	// BreakStmt
	`0x7fca2d8070e0 <col:11, col:23>`:
	ast.BreakStmt{
		Address: "0x7fca2d8070e0",
		Position: "col:11, col:23",
	},

	// BuiltinType
	`0x7f8a43023f40 '__int128'`:
	ast.BuiltinType{
		Address: "0x7f8a43023f40",
		Type: "__int128",
	},
	`0x7f8a43023ea0 'unsigned long long'`:
	ast.BuiltinType{
		Address: "0x7f8a43023ea0",
		Type: "unsigned long long",
	},

	// CallExpr
	`0x7f9bf3033240 <col:11, col:25> 'int'`:
	ast.CallExpr{
		Address: "0x7f9bf3033240",
		Position: "col:11, col:25",
		Type: "int",
	},
	`0x7f9bf3035c20 <line:7:4, col:64> 'int'`:
	ast.CallExpr{
		Address: "0x7f9bf3035c20",
		Position: "line:7:4, col:64",
		Type: "int",
	},

	// CharacterLiteral
	`0x7f980b858308 <col:62> 'int' 10`:
	ast.CharacterLiteral{
		Address: "0x7f980b858308",
		Position: "col:62",
		Type: "int",
		Value: 10,
	},

	// CompoundStmt
	`0x7fbd0f014f18 <col:54, line:358:1>`:
	ast.CompoundStmt{
		Address: "0x7fbd0f014f18",
		Position: "col:54, line:358:1",
	},
	`0x7fbd0f8360b8 <line:4:1, line:13:1>`:
	ast.CompoundStmt{
		Address: "0x7fbd0f8360b8",
		Position: "line:4:1, line:13:1",
	},

	// ConstantArrayType
	`0x7f94ad016a40 'struct __va_list_tag [1]' 1`:
	ast.ConstantArrayType{
		Address: "0x7f94ad016a40",
		Type: "struct __va_list_tag [1]",
		Size: 1,
	},
	`0x7f8c5f059d20 'char [37]' 37`:
	ast.ConstantArrayType{
		Address: "0x7f8c5f059d20",
		Type: "char [37]",
		Size: 37,
	},

	// CStyleCastExpr
	`0x7fddc18fb2e0 <col:50, col:56> 'char' <IntegralCast>`:
	ast.CStyleCastExpr{
		Address: "0x7fddc18fb2e0",
		Position: "col:50, col:56",
		Type: "char",
		Kind: "IntegralCast",
	},

	// DeclRefExpr
	`0x7fc972064460 <col:8> 'FILE *' lvalue ParmVar 0x7fc9720642d0 '_p' 'FILE *'`:
	ast.DeclRefExpr{
		Address: "0x7fc972064460",
		Position: "col:8",
		Type: "FILE *",
		Lvalue: true,
		For: "ParmVar",
		Address2: "0x7fc9720642d0",
		Name: "_p",
		Type2: "FILE *",
	},
	`0x7fc97206a958 <col:11> 'int (int, FILE *)' Function 0x7fc972064198 '__swbuf' 'int (int, FILE *)'`:
	ast.DeclRefExpr{
		Address: "0x7fc97206a958",
		Position: "col:11",
		Type: "int (int, FILE *)",
		Lvalue: false,
		For: "Function",
		Address2: "0x7fc972064198",
		Name: "__swbuf",
		Type2: "int (int, FILE *)",
	},
	`0x7fa36680f170 <col:19> 'struct programming':'struct programming' lvalue Var 0x7fa36680dc20 'variable' 'struct programming':'struct programming'`:
	ast.DeclRefExpr{
		Address: "0x7fa36680f170",
		Position: "col:19",
		Type: "struct programming",
		Lvalue: true,
		For: "Var",
		Address2: "0x7fa36680dc20",
		Name: "variable",
		Type2: "struct programming",
	},

	// DeclStmt
	`0x7fb791846e80 <line:11:4, col:31>`:
	ast.DeclStmt{
		Address: "0x7fb791846e80",
		Position: "line:11:4, col:31",
	},

	// DeprecatedAttr
	`0x7fec4b0ab9c0 <line:180:48, col:63> "This function is provided for compatibility reasons only.  Due to security concerns inherent in the design of tempnam(3), it is highly recommended that you use mkstemp(3) instead." ""`:
	ast.DeprecatedAttr{
		Address: "0x7fec4b0ab9c0",
		Position: "line:180:48, col:63",
		Message1: "This function is provided for compatibility reasons only.  Due to security concerns inherent in the design of tempnam(3), it is highly recommended that you use mkstemp(3) instead.",
		Message2: "",
	},

	// ElaboratedType
	`0x7f873686c120 'union __mbstate_t' sugar`:
	ast.ElaboratedType{
		Address: "0x7f873686c120",
		Type: "union __mbstate_t",
		Tags: "sugar",
	},

	// FieldDecl
	`0x7fef510c4848 <line:141:2, col:6> col:6 _ur 'int'`:
	ast.FieldDecl{
		Address: "0x7fef510c4848",
		Position: "line:141:2, col:6",
		Position2: "col:6",
		Name: "_ur",
		Type: "int",
		Referenced: false,
	},
	`0x7fef510c46f8 <line:139:2, col:16> col:16 _ub 'struct __sbuf':'struct __sbuf'`:
	ast.FieldDecl{
		Address: "0x7fef510c46f8",
		Position: "line:139:2, col:16",
		Position2: "col:16",
		Name: "_ub",
		Type: "struct __sbuf",
		Referenced: false,
	},
	`0x7fef510c3fe0 <line:134:2, col:19> col:19 _read 'int (* _Nullable)(void *, char *, int)':'int (*)(void *, char *, int)'`:
	ast.FieldDecl{
		Address: "0x7fef510c3fe0",
		Position: "line:134:2, col:19",
		Position2: "col:19",
		Name: "_read",
		Type: "int (* _Nullable)(void *, char *, int)",
		Referenced: false,
	},
	`0x7fef51073a60 <line:105:2, col:40> col:40 __cleanup_stack 'struct __darwin_pthread_handler_rec *'`:
	ast.FieldDecl{
		Address: "0x7fef51073a60",
		Position: "line:105:2, col:40",
		Position2: "col:40",
		Name: "__cleanup_stack",
		Type: "struct __darwin_pthread_handler_rec *",
		Referenced: false,
	},
	`0x7fef510738e8 <line:100:2, col:43> col:7 __opaque 'char [16]'`:
	ast.FieldDecl{
		Address: "0x7fef510738e8",
		Position: "line:100:2, col:43",
		Position2: "col:7",
		Name: "__opaque",
		Type: "char [16]",
		Referenced: false,
	},
	`0x7fe9f5072268 <line:129:2, col:6> col:6 referenced _lbfsize 'int'`:
	ast.FieldDecl{
		Address: "0x7fe9f5072268",
		Position: "line:129:2, col:6",
		Position2: "col:6",
		Name: "_lbfsize",
		Type: "int",
		Referenced: true,
	},
}

func TestNodes(t *testing.T) {
	for line, expected := range nodes {
		// Append the name of the struct onto the front. This would make
		// the complete line it would normally be parsing.
		actual := ast.Parse(reflect.TypeOf(expected).Name() + " " + line)

		if expected != actual {
			t.Errorf("\nexpected: %#v\n     got: %#v", expected, actual)
		}
	}
}
