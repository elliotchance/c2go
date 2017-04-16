package main

import (
	"reflect"
	"testing"
)

var nodes = map[string]interface{}{
	// AlwaysInlineAttr
	`0x7fce780f5018 </usr/include/sys/cdefs.h:313:68> always_inline`: &AlwaysInlineAttr{
		Address:  "0x7fce780f5018",
		Position: "/usr/include/sys/cdefs.h:313:68",
		Children: []interface{}{},
	},

	// ArraySubscriptExpr
	`0x7fe35b85d180 <col:63, col:69> 'char *' lvalue`: &ArraySubscriptExpr{
		Address:  "0x7fe35b85d180",
		Position: "col:63, col:69",
		Type:     "char *",
		Kind:     "lvalue",
		Children: []interface{}{},
	},

	// AsmLabelAttr
	`0x7ff26d8224e8 </usr/include/sys/cdefs.h:569:36> "_fopen"`: &AsmLabelAttr{
		Address:      "0x7ff26d8224e8",
		Position:     "/usr/include/sys/cdefs.h:569:36",
		FunctionName: "_fopen",
		Children:     []interface{}{},
	},

	// AvailabilityAttr
	`0x7fc5ff8e5d18 </usr/include/AvailabilityInternal.h:21697:88, col:124> macos 10.10 0 0 "" ""`: &AvailabilityAttr{
		Address:     "0x7fc5ff8e5d18",
		Position:    "/usr/include/AvailabilityInternal.h:21697:88, col:124",
		OS:          "macos",
		Version:     "10.10",
		Unknown1:    0,
		Unknown2:    0,
		Unavailable: false,
		Message1:    "",
		Message2:    "",
		Children:    []interface{}{},
	},
	`0x7fc5ff8e60d0 </usr/include/Availability.h:215:81, col:115> watchos 3.0 0 0 "" ""`: &AvailabilityAttr{
		Address:     "0x7fc5ff8e60d0",
		Position:    "/usr/include/Availability.h:215:81, col:115",
		OS:          "watchos",
		Version:     "3.0",
		Unknown1:    0,
		Unknown2:    0,
		Unavailable: false,
		Message1:    "",
		Message2:    "",
		Children:    []interface{}{},
	},
	`0x7fc5ff8e6170 <col:81, col:115> tvos 10.0 0 0 "" ""`: &AvailabilityAttr{
		Address:     "0x7fc5ff8e6170",
		Position:    "col:81, col:115",
		OS:          "tvos",
		Version:     "10.0",
		Unknown1:    0,
		Unknown2:    0,
		Unavailable: false,
		Message1:    "",
		Message2:    "",
		Children:    []interface{}{},
	},
	`0x7fc5ff8e61d8 <col:81, col:115> ios 10.0 0 0 "" ""`: &AvailabilityAttr{
		Address:     "0x7fc5ff8e61d8",
		Position:    "col:81, col:115",
		OS:          "ios",
		Version:     "10.0",
		Unknown1:    0,
		Unknown2:    0,
		Unavailable: false,
		Message1:    "",
		Message2:    "",
		Children:    []interface{}{},
	},
	`0x7fc5ff8f0e18 </usr/include/sys/cdefs.h:275:50, col:99> swift 0 0 0 Unavailable "Use snprintf instead." ""`: &AvailabilityAttr{
		Address:     "0x7fc5ff8f0e18",
		Position:    "/usr/include/sys/cdefs.h:275:50, col:99",
		OS:          "swift",
		Version:     "0",
		Unknown1:    0,
		Unknown2:    0,
		Unavailable: true,
		Message1:    "Use snprintf instead.",
		Message2:    "",
		Children:    []interface{}{},
	},
	`0x7fc5ff8f1988 <line:275:50, col:99> swift 0 0 0 Unavailable "Use mkstemp(3) instead." ""`: &AvailabilityAttr{
		Address:     "0x7fc5ff8f1988",
		Position:    "line:275:50, col:99",
		OS:          "swift",
		Version:     "0",
		Unknown1:    0,
		Unknown2:    0,
		Unavailable: true,
		Message1:    "Use mkstemp(3) instead.",
		Message2:    "",
		Children:    []interface{}{},
	},
	`0x104035438 </usr/include/AvailabilityInternal.h:14571:88, col:124> macosx 10.10 0 0 ""`: &AvailabilityAttr{
		Address:     "0x104035438",
		Position:    "/usr/include/AvailabilityInternal.h:14571:88, col:124",
		OS:          "macosx",
		Version:     "10.10",
		Unknown1:    0,
		Unknown2:    0,
		Unavailable: false,
		Message1:    "",
		Message2:    "",
		Children:    []interface{}{},
	},

	// BinaryOperator
	`0x7fca2d8070e0 <col:11, col:23> 'unsigned char' '='`: &BinaryOperator{
		Address:  "0x7fca2d8070e0",
		Position: "col:11, col:23",
		Type:     "unsigned char",
		Operator: "=",
		Children: []interface{}{},
	},

	// BreakStmt
	`0x7fca2d8070e0 <col:11, col:23>`: &BreakStmt{
		Address:  "0x7fca2d8070e0",
		Position: "col:11, col:23",
		Children: []interface{}{},
	},

	// BuiltinType
	`0x7f8a43023f40 '__int128'`: &BuiltinType{
		Address:  "0x7f8a43023f40",
		Type:     "__int128",
		Children: []interface{}{},
	},
	`0x7f8a43023ea0 'unsigned long long'`: &BuiltinType{
		Address:  "0x7f8a43023ea0",
		Type:     "unsigned long long",
		Children: []interface{}{},
	},

	// CallExpr
	`0x7f9bf3033240 <col:11, col:25> 'int'`: &CallExpr{
		Address:  "0x7f9bf3033240",
		Position: "col:11, col:25",
		Type:     "int",
		Children: []interface{}{},
	},
	`0x7f9bf3035c20 <line:7:4, col:64> 'int'`: &CallExpr{
		Address:  "0x7f9bf3035c20",
		Position: "line:7:4, col:64",
		Type:     "int",
		Children: []interface{}{},
	},

	// CharacterLiteral
	`0x7f980b858308 <col:62> 'int' 10`: &CharacterLiteral{
		Address:  "0x7f980b858308",
		Position: "col:62",
		Type:     "int",
		Value:    10,
		Children: []interface{}{},
	},

	// CompoundStmt
	`0x7fbd0f014f18 <col:54, line:358:1>`: &CompoundStmt{
		Address:  "0x7fbd0f014f18",
		Position: "col:54, line:358:1",
		Children: []interface{}{},
	},
	`0x7fbd0f8360b8 <line:4:1, line:13:1>`: &CompoundStmt{
		Address:  "0x7fbd0f8360b8",
		Position: "line:4:1, line:13:1",
		Children: []interface{}{},
	},

	// ConditionalOperator
	`0x7fc6ae0bc678 <col:6, col:89> 'void'`: &ConditionalOperator{
		Address:  "0x7fc6ae0bc678",
		Position: "col:6, col:89",
		Type:     "void",
		Children: []interface{}{},
	},

	// ConstAttr
	`0x7fa3b88bbb38 <line:4:1, line:13:1>foo`: &ConstAttr{
		Address:  "0x7fa3b88bbb38",
		Position: "line:4:1, line:13:1",
		Tags:     "foo",
		Children: []interface{}{},
	},

	// ConstantArrayType
	`0x7f94ad016a40 'struct __va_list_tag [1]' 1`: &ConstantArrayType{
		Address:  "0x7f94ad016a40",
		Type:     "struct __va_list_tag [1]",
		Size:     1,
		Children: []interface{}{},
	},
	`0x7f8c5f059d20 'char [37]' 37`: &ConstantArrayType{
		Address:  "0x7f8c5f059d20",
		Type:     "char [37]",
		Size:     37,
		Children: []interface{}{},
	},

	// CStyleCastExpr
	`0x7fddc18fb2e0 <col:50, col:56> 'char' <IntegralCast>`: &CStyleCastExpr{
		Address:  "0x7fddc18fb2e0",
		Position: "col:50, col:56",
		Type:     "char",
		Kind:     "IntegralCast",
		Children: []interface{}{},
	},

	// DeclRefExpr
	`0x7fc972064460 <col:8> 'FILE *' lvalue ParmVar 0x7fc9720642d0 '_p' 'FILE *'`: &DeclRefExpr{
		Address:  "0x7fc972064460",
		Position: "col:8",
		Type:     "FILE *",
		Lvalue:   true,
		For:      "ParmVar",
		Address2: "0x7fc9720642d0",
		Name:     "_p",
		Type2:    "FILE *",
		Children: []interface{}{},
	},
	`0x7fc97206a958 <col:11> 'int (int, FILE *)' Function 0x7fc972064198 '__swbuf' 'int (int, FILE *)'`: &DeclRefExpr{
		Address:  "0x7fc97206a958",
		Position: "col:11",
		Type:     "int (int, FILE *)",
		Lvalue:   false,
		For:      "Function",
		Address2: "0x7fc972064198",
		Name:     "__swbuf",
		Type2:    "int (int, FILE *)",
		Children: []interface{}{},
	},
	`0x7fa36680f170 <col:19> 'struct programming':'struct programming' lvalue Var 0x7fa36680dc20 'variable' 'struct programming':'struct programming'`: &DeclRefExpr{
		Address:  "0x7fa36680f170",
		Position: "col:19",
		Type:     "struct programming",
		Lvalue:   true,
		For:      "Var",
		Address2: "0x7fa36680dc20",
		Name:     "variable",
		Type2:    "struct programming",
		Children: []interface{}{},
	},

	// DeclStmt
	`0x7fb791846e80 <line:11:4, col:31>`: &DeclStmt{
		Address:  "0x7fb791846e80",
		Position: "line:11:4, col:31",
		Children: []interface{}{},
	},

	// DeprecatedAttr
	`0x7fec4b0ab9c0 <line:180:48, col:63> "This function is provided for compatibility reasons only.  Due to security concerns inherent in the design of tempnam(3), it is highly recommended that you use mkstemp(3) instead." ""`: &DeprecatedAttr{
		Address:  "0x7fec4b0ab9c0",
		Position: "line:180:48, col:63",
		Message1: "This function is provided for compatibility reasons only.  Due to security concerns inherent in the design of tempnam(3), it is highly recommended that you use mkstemp(3) instead.",
		Message2: "",
		Children: []interface{}{},
	},

	// ElaboratedType
	`0x7f873686c120 'union __mbstate_t' sugar`: &ElaboratedType{
		Address:  "0x7f873686c120",
		Type:     "union __mbstate_t",
		Tags:     "sugar",
		Children: []interface{}{},
	},

	// Enum
	`0x7f980b858308 'foo'`: &Enum{
		Address:  "0x7f980b858308",
		Name:     "foo",
		Children: []interface{}{},
	},

	// EnumDecl
	`0x22a6c80 <line:180:1, line:186:1> __codecvt_result`: &EnumDecl{
		Address:   "0x22a6c80",
		Position:  "line:180:1, line:186:1",
		Position2: "",
		Name:      "__codecvt_result",
		Children:  []interface{}{},
	},

	// EnumConstantDecl
	`0x1660db0 <line:185:3> __codecvt_noconv 'int'`: &EnumConstantDecl{
		Address:   "0x1660db0",
		Position:  "line:185:3",
		Position2: "",
		Name:      "__codecvt_noconv",
		Type:      "int",
		Children:  []interface{}{},
	},

	// EnumType
	`0x7f980b858309 'foo'`: &EnumType{
		Address:  "0x7f980b858309",
		Name:     "foo",
		Children: []interface{}{},
	},

	// FieldDecl
	`0x7fef510c4848 <line:141:2, col:6> col:6 _ur 'int'`: &FieldDecl{
		Address:    "0x7fef510c4848",
		Position:   "line:141:2, col:6",
		Position2:  "col:6",
		Name:       "_ur",
		Type:       "int",
		Referenced: false,
		Children:   []interface{}{},
	},
	`0x7fef510c46f8 <line:139:2, col:16> col:16 _ub 'struct __sbuf':'struct __sbuf'`: &FieldDecl{
		Address:    "0x7fef510c46f8",
		Position:   "line:139:2, col:16",
		Position2:  "col:16",
		Name:       "_ub",
		Type:       "struct __sbuf",
		Referenced: false,
		Children:   []interface{}{},
	},
	`0x7fef510c3fe0 <line:134:2, col:19> col:19 _read 'int (* _Nullable)(void *, char *, int)':'int (*)(void *, char *, int)'`: &FieldDecl{
		Address:    "0x7fef510c3fe0",
		Position:   "line:134:2, col:19",
		Position2:  "col:19",
		Name:       "_read",
		Type:       "int (* _Nullable)(void *, char *, int)",
		Referenced: false,
		Children:   []interface{}{},
	},
	`0x7fef51073a60 <line:105:2, col:40> col:40 __cleanup_stack 'struct __darwin_pthread_handler_rec *'`: &FieldDecl{
		Address:    "0x7fef51073a60",
		Position:   "line:105:2, col:40",
		Position2:  "col:40",
		Name:       "__cleanup_stack",
		Type:       "struct __darwin_pthread_handler_rec *",
		Referenced: false,
		Children:   []interface{}{},
	},
	`0x7fef510738e8 <line:100:2, col:43> col:7 __opaque 'char [16]'`: &FieldDecl{
		Address:    "0x7fef510738e8",
		Position:   "line:100:2, col:43",
		Position2:  "col:7",
		Name:       "__opaque",
		Type:       "char [16]",
		Referenced: false,
		Children:   []interface{}{},
	},
	`0x7fe9f5072268 <line:129:2, col:6> col:6 referenced _lbfsize 'int'`: &FieldDecl{
		Address:    "0x7fe9f5072268",
		Position:   "line:129:2, col:6",
		Position2:  "col:6",
		Name:       "_lbfsize",
		Type:       "int",
		Referenced: true,
		Children:   []interface{}{},
	},
	`0x7f9bc9083d00 <line:91:5, line:97:8> line:91:5 'unsigned short'`: &FieldDecl{
		Address:    "0x7f9bc9083d00",
		Position:   "line:91:5, line:97:8",
		Position2:  "line:91:5",
		Name:       "",
		Type:       "unsigned short",
		Referenced: false,
		Children:   []interface{}{},
	},
	`0x30363a0 <col:18, col:29> __val 'int [2]'`: &FieldDecl{
		Address:    "0x30363a0",
		Position:   "col:18, col:29",
		Position2:  "",
		Name:       "__val",
		Type:       "int [2]",
		Referenced: false,
		Children:   []interface{}{},
	},

	// FloatingLiteral
	`0x7febe106f5e8 <col:24> 'double' 1.230000e+00`: &FloatingLiteral{
		Address:  "0x7febe106f5e8",
		Position: "col:24",
		Type:     "double",
		Value:    1.23,
		Children: []interface{}{},
	},

	// FormatAttr
	`0x7fcc8d8ecee8 <col:6> Implicit printf 2 3`: &FormatAttr{
		Address:      "0x7fcc8d8ecee8",
		Position:     "col:6",
		Implicit:     true,
		Inherited:    false,
		FunctionName: "printf",
		Unknown1:     2,
		Unknown2:     3,
		Children:     []interface{}{},
	},
	`0x7fcc8d8ecff8 </usr/include/sys/cdefs.h:351:18, col:61> printf 2 3`: &FormatAttr{
		Address:      "0x7fcc8d8ecff8",
		Position:     "/usr/include/sys/cdefs.h:351:18, col:61",
		Implicit:     false,
		Inherited:    false,
		FunctionName: "printf",
		Unknown1:     2,
		Unknown2:     3,
		Children:     []interface{}{},
	},
	`0x273b4d0 <line:357:12> Inherited printf 2 3`: &FormatAttr{
		Address:      "0x273b4d0",
		Position:     "line:357:12",
		Implicit:     false,
		Inherited:    true,
		FunctionName: "printf",
		Unknown1:     2,
		Unknown2:     3,
		Children:     []interface{}{},
	},

	// FunctionDecl
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
		Children:   []interface{}{},
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
		Children:   []interface{}{},
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
		Children:   []interface{}{},
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
		Children:   []interface{}{},
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
		Children:   []interface{}{},
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
		Children:   []interface{}{},
	},

	// FunctionProtoType
	`0x7fa3b88bbb30 'struct _opaque_pthread_t *' foo`: &FunctionProtoType{
		Address:  "0x7fa3b88bbb30",
		Type:     "struct _opaque_pthread_t *",
		Kind:     "foo",
		Children: []interface{}{},
	},

	// ForStmt
	`0x7f961e018848 <line:9:4, line:10:70>`: &ForStmt{
		Address:  "0x7f961e018848",
		Position: "line:9:4, line:10:70",
		Children: []interface{}{},
	},

	// IfStmt
	`0x7fc0a69091d0 <line:11:7, line:18:7>`: &IfStmt{
		Address:  "0x7fc0a69091d0",
		Position: "line:11:7, line:18:7",
		Children: []interface{}{},
	},

	// ImplicitCastExpr
	`0x7f9f5b0a1288 <col:8> 'FILE *' <LValueToRValue>`: &ImplicitCastExpr{
		Address:  "0x7f9f5b0a1288",
		Position: "col:8",
		Type:     "FILE *",
		Kind:     "LValueToRValue",
		Children: []interface{}{},
	},
	`0x7f9f5b0a7828 <col:11> 'int (*)(int, FILE *)' <FunctionToPointerDecay>`: &ImplicitCastExpr{
		Address:  "0x7f9f5b0a7828",
		Position: "col:11",
		Type:     "int (*)(int, FILE *)",
		Kind:     "FunctionToPointerDecay",
		Children: []interface{}{},
	},

	// IntegerLiteral
	`0x7fbe9804bcc8 <col:14> 'int' 1`: &IntegerLiteral{
		Address:  "0x7fbe9804bcc8",
		Position: "col:14",
		Type:     "int",
		Value:    1,
		Children: []interface{}{},
	},

	// MallocAttr
	`0x7fc0a69091d1 <line:11:7, line:18:7>`: &MallocAttr{
		Address:  "0x7fc0a69091d1",
		Position: "line:11:7, line:18:7",
		Children: []interface{}{},
	},

	// MemberExpr
	`0x7fcc758e34a0 <col:8, col:12> 'int' lvalue ->_w 0x7fcc758d60c8`: &MemberExpr{
		Address:  "0x7fcc758e34a0",
		Position: "col:8, col:12",
		Type:     "int",
		Lvalue:   true,
		Name:     "_w",
		Address2: "0x7fcc758d60c8",
		Children: []interface{}{},
	},
	`0x7fcc76004210 <col:12, col:16> 'unsigned char *' lvalue ->_p 0x7fcc758d6018`: &MemberExpr{
		Address:  "0x7fcc76004210",
		Position: "col:12, col:16",
		Type:     "unsigned char *",
		Lvalue:   true,
		Name:     "_p",
		Address2: "0x7fcc758d6018",
		Children: []interface{}{},
	},
	`0x7f85338325b0 <col:4, col:13> 'float' lvalue .constant 0x7f8533832260`: &MemberExpr{
		Address:  "0x7f85338325b0",
		Position: "col:4, col:13",
		Type:     "float",
		Lvalue:   true,
		Name:     "constant",
		Address2: "0x7f8533832260",
		Children: []interface{}{},
	},
	`0x7f8533832670 <col:4, col:13> 'char *' lvalue .pointer 0x7f85338322b8`: &MemberExpr{
		Address:  "0x7f8533832670",
		Position: "col:4, col:13",
		Type:     "char *",
		Lvalue:   true,
		Name:     "pointer",
		Address2: "0x7f85338322b8",
		Children: []interface{}{},
	},

	// ModeAttr
	`0x7f980b858309 <line:11:7, line:18:7> foo`: &ModeAttr{
		Address:  "0x7f980b858309",
		Position: "line:11:7, line:18:7",
		Name:     "foo",
		Children: []interface{}{},
	},

	// NoThrowAttr
	`0x7fa1488273a0 <line:7:4, line:11:4>`: &NoThrowAttr{
		Address:  "0x7fa1488273a0",
		Position: "line:7:4, line:11:4",
		Children: []interface{}{},
	},

	// NonNullAttr
	`0x7fa1488273b0 <line:7:4, line:11:4> 1`: &NonNullAttr{
		Address:  "0x7fa1488273b0",
		Position: "line:7:4, line:11:4",
		Children: []interface{}{},
	},
	`0x2cce280 </sys/cdefs.h:286:44, /bits/mathcalls.h:115:69> 1`: &NonNullAttr{
		Address:  "0x2cce280",
		Position: "/sys/cdefs.h:286:44, /bits/mathcalls.h:115:69",
		Children: []interface{}{},
	},

	// ParenExpr
	`0x7fb0bc8b2308 <col:10, col:25> 'unsigned char'`: &ParenExpr{
		Address:  "0x7fb0bc8b2308",
		Position: "col:10, col:25",
		Type:     "unsigned char",
		Children: []interface{}{},
	},

	// ParmVarDecl
	`0x7f973380f000 <col:14> col:17 'int'`: &ParmVarDecl{
		Address:   "0x7f973380f000",
		Position:  "col:14",
		Position2: "col:17",
		Type:      "int",
		Name:      "",
		Type2:     "",
		IsUsed:    false,
		Children:  []interface{}{},
	},
	`0x7f973380f070 <col:19, col:30> col:31 'const char *'`: &ParmVarDecl{
		Address:   "0x7f973380f070",
		Position:  "col:19, col:30",
		Position2: "col:31",
		Type:      "const char *",
		Name:      "",
		Type2:     "",
		IsUsed:    false,
		Children:  []interface{}{},
	},
	`0x7f9733816e50 <col:13, col:37> col:37 __filename 'const char *__restrict'`: &ParmVarDecl{
		Address:   "0x7f9733816e50",
		Position:  "col:13, col:37",
		Position2: "col:37",
		Type:      "const char *__restrict",
		Name:      "__filename",
		Type2:     "",
		IsUsed:    false,
		Children:  []interface{}{},
	},
	`0x7f9733817418 <<invalid sloc>> <invalid sloc> 'FILE *'`: &ParmVarDecl{
		Address:   "0x7f9733817418",
		Position:  "<invalid sloc>",
		Position2: "<invalid sloc>",
		Type:      "FILE *",
		Name:      "",
		Type2:     "",
		IsUsed:    false,
		Children:  []interface{}{},
	},
	`0x7f9733817c30 <col:40, col:47> col:47 __size 'size_t':'unsigned long'`: &ParmVarDecl{
		Address:   "0x7f9733817c30",
		Position:  "col:40, col:47",
		Position2: "col:47",
		Type:      "size_t",
		Name:      "__size",
		Type2:     "unsigned long",
		IsUsed:    false,
		Children:  []interface{}{},
	},
	`0x7f973382fa10 <line:476:18, col:25> col:34 'int (* _Nullable)(void *, char *, int)':'int (*)(void *, char *, int)'`: &ParmVarDecl{
		Address:   "0x7f973382fa10",
		Position:  "line:476:18, col:25",
		Position2: "col:34",
		Type:      "int (* _Nullable)(void *, char *, int)",
		Name:      "",
		Type2:     "int (*)(void *, char *, int)",
		IsUsed:    false,
		Children:  []interface{}{},
	},
	`0x7f97338355b8 <col:10, col:14> col:14 used argc 'int'`: &ParmVarDecl{
		Address:   "0x7f97338355b8",
		Position:  "col:10, col:14",
		Position2: "col:14",
		Type:      "int",
		Name:      "argc",
		Type2:     "",
		IsUsed:    true,
		Children:  []interface{}{},
	},

	// PointerType
	`0x7fa3b88bbb30 'struct _opaque_pthread_t *'`: &PointerType{
		Address:  "0x7fa3b88bbb30",
		Type:     "struct _opaque_pthread_t *",
		Children: []interface{}{},
	},

	// PredefinedExpr
	`0x33d6e08 <col:30> 'const char [25]' lvalue __PRETTY_FUNCTION__`: &PredefinedExpr{
		Address:  "0x33d6e08",
		Position: "col:30",
		Type:     "const char [25]",
		Lvalue:   true,
		Name:     "__PRETTY_FUNCTION__",
		Children: []interface{}{},
	},

	// QualType
	`0x7fa3b88bbb31 'struct _opaque_pthread_t *' foo`: &QualType{
		Address:  "0x7fa3b88bbb31",
		Type:     "struct _opaque_pthread_t *",
		Kind:     "foo",
		Children: []interface{}{},
	},

	// Record
	`0x7fd3ab857950 '__sFILE'`: &Record{
		Address:  "0x7fd3ab857950",
		Type:     "__sFILE",
		Children: []interface{}{},
	},

	// RecordDecl
	`0x7f913c0dbb50 <line:76:9, line:79:1> line:76:9 union definition`: &RecordDecl{
		Address:    "0x7f913c0dbb50",
		Position:   "line:76:9, line:79:1",
		Prev:       "",
		Position2:  "line:76:9",
		Kind:       "union",
		Name:       "",
		Definition: true,
		Children:   []interface{}{},
	},
	`0x7f85360285c8 </usr/include/sys/_pthread/_pthread_types.h:57:1, line:61:1> line:57:8 struct __darwin_pthread_handler_rec definition`: &RecordDecl{
		Address:    "0x7f85360285c8",
		Position:   "/usr/include/sys/_pthread/_pthread_types.h:57:1, line:61:1",
		Prev:       "",
		Position2:  "line:57:8",
		Kind:       "struct",
		Name:       "__darwin_pthread_handler_rec",
		Definition: true,
		Children:   []interface{}{},
	},
	`0x7f85370248a0 <line:94:1, col:8> col:8 struct __sFILEX`: &RecordDecl{
		Address:    "0x7f85370248a0",
		Position:   "line:94:1, col:8",
		Prev:       "",
		Position2:  "col:8",
		Kind:       "struct",
		Name:       "__sFILEX",
		Definition: false,
		Children:   []interface{}{},
	},

	// RecordType
	`0x7fd3ab84dda0 'struct _opaque_pthread_condattr_t'`: &RecordType{
		Address:  "0x7fd3ab84dda0",
		Type:     "struct _opaque_pthread_condattr_t",
		Children: []interface{}{},
	},

	// RestrictAttr
	`0x7f980b858305 <line:11:7, line:18:7> foo`: &RestrictAttr{
		Address:  "0x7f980b858305",
		Position: "line:11:7, line:18:7",
		Name:     "foo",
		Children: []interface{}{},
	},

	// ReturnStmt
	`0x7fbb7a8325e0 <line:13:4, col:11>`: &ReturnStmt{
		Address:  "0x7fbb7a8325e0",
		Position: "line:13:4, col:11",
		Children: []interface{}{},
	},

	// StringLiteral
	`0x7fe16f0b4d58 <col:11> 'char [45]' lvalue "Number of command line arguments passed: %d\n"`: &StringLiteral{
		Address:  "0x7fe16f0b4d58",
		Position: "col:11",
		Type:     "char [45]",
		Lvalue:   true,
		Value:    "Number of command line arguments passed: %d\n",
		Children: []interface{}{},
	},

	// TranslationUnitDecl
	`0x7fe78a815ed0 <<invalid sloc>> <invalid sloc>`: &TranslationUnitDecl{
		Address:  "0x7fe78a815ed0",
		Children: []interface{}{},
	},

	// Typedef
	`0x7f84d10dc1d0 '__darwin_ssize_t'`: &Typedef{
		Address:  "0x7f84d10dc1d0",
		Type:     "__darwin_ssize_t",
		Children: []interface{}{},
	},

	// TypedefDecl
	`0x7fdef0862430 <line:120:1, col:16> col:16`: &TypedefDecl{
		Address:      "0x7fdef0862430",
		Position:     "line:120:1, col:16",
		Position2:    "col:16",
		Name:         "",
		Type:         "",
		Type2:        "",
		IsImplicit:   false,
		IsReferenced: false,
		Children:     []interface{}{},
	},
	`0x7ffb9f824278 <<invalid sloc>> <invalid sloc> implicit __uint128_t 'unsigned __int128'`: &TypedefDecl{
		Address:      "0x7ffb9f824278",
		Position:     "<invalid sloc>",
		Position2:    "<invalid sloc>",
		Name:         "__uint128_t",
		Type:         "unsigned __int128",
		Type2:        "",
		IsImplicit:   true,
		IsReferenced: false,
		Children:     []interface{}{},
	},
	`0x7ffb9f824898 <<invalid sloc>> <invalid sloc> implicit referenced __builtin_va_list 'struct __va_list_tag [1]'`: &TypedefDecl{
		Address:      "0x7ffb9f824898",
		Position:     "<invalid sloc>",
		Position2:    "<invalid sloc>",
		Name:         "__builtin_va_list",
		Type:         "struct __va_list_tag [1]",
		Type2:        "",
		IsImplicit:   true,
		IsReferenced: true,
		Children:     []interface{}{},
	},
	`0x7ffb9f8248f8 </usr/include/i386/_types.h:37:1, col:24> col:24 __int8_t 'signed char'`: &TypedefDecl{
		Address:      "0x7ffb9f8248f8",
		Position:     "/usr/include/i386/_types.h:37:1, col:24",
		Position2:    "col:24",
		Name:         "__int8_t",
		Type:         "signed char",
		Type2:        "",
		IsImplicit:   false,
		IsReferenced: false,
		Children:     []interface{}{},
	},
	`0x7ffb9f8dbf50 <line:98:1, col:27> col:27 referenced __darwin_va_list '__builtin_va_list':'struct __va_list_tag [1]'`: &TypedefDecl{
		Address:      "0x7ffb9f8dbf50",
		Position:     "line:98:1, col:27",
		Position2:    "col:27",
		Name:         "__darwin_va_list",
		Type:         "__builtin_va_list",
		Type2:        "struct __va_list_tag [1]",
		IsImplicit:   false,
		IsReferenced: true,
		Children:     []interface{}{},
	},
	`0x34461f0 <line:338:1, col:77> __io_read_fn '__ssize_t (void *, char *, size_t)'`: &TypedefDecl{
		Address:      "0x34461f0",
		Position:     "line:338:1, col:77",
		Position2:    "",
		Name:         "__io_read_fn",
		Type:         "__ssize_t (void *, char *, size_t)",
		Type2:        "",
		IsImplicit:   false,
		IsReferenced: false,
		Children:     []interface{}{},
	},
	// Issue: #26
	`0x55b9da8784b0 <line:341:1, line:342:16> line:341:19 __io_write_fn '__ssize_t (void *, const char *, size_t)'`: &TypedefDecl{
		Address:      "0x55b9da8784b0",
		Position:     "line:341:1, line:342:16",
		Position2:    "line:341:19",
		Name:         "__io_write_fn",
		Type:         "__ssize_t (void *, const char *, size_t)",
		Type2:        "",
		IsImplicit:   false,
		IsReferenced: false,
		Children:     []interface{}{},
	},

	// TypedefType
	`0x7f887a0dc760 '__uint16_t' sugar`: &TypedefType{
		Address:  "0x7f887a0dc760",
		Type:     "__uint16_t",
		Tags:     "sugar",
		Children: []interface{}{},
	},

	// UnaryOperator
	`0x7fe0260f50d8 <col:6, col:12> 'int' prefix '--'`: &UnaryOperator{
		Address:  "0x7fe0260f50d8",
		Position: "col:6, col:12",
		Type:     "int",
		IsLvalue: false,
		IsPrefix: true,
		Operator: "--",
		Children: []interface{}{},
	},
	`0x7fe0260fb468 <col:11, col:18> 'unsigned char' lvalue prefix '*'`: &UnaryOperator{
		Address:  "0x7fe0260fb468",
		Position: "col:11, col:18",
		Type:     "unsigned char",
		IsLvalue: true,
		IsPrefix: true,
		Operator: "*",
		Children: []interface{}{},
	},
	`0x7fe0260fb448 <col:12, col:18> 'unsigned char *' postfix '++'`: &UnaryOperator{
		Address:  "0x7fe0260fb448",
		Position: "col:12, col:18",
		Type:     "unsigned char *",
		IsLvalue: false,
		IsPrefix: false,
		Operator: "++",
		Children: []interface{}{},
	},

	// VarDecl
	`0x7fd5e90e5a00 <col:14> col:17 'int'`: &VarDecl{
		Address:   "0x7fd5e90e5a00",
		Position:  "col:14",
		Position2: "col:17",
		Name:      "",
		Type:      "int",
		Type2:     "",
		IsExtern:  false,
		IsUsed:    false,
		IsCInit:   false,
		Children:  []interface{}{},
	},
	`0x7fd5e90e9078 <line:156:1, col:14> col:14 __stdinp 'FILE *' extern`: &VarDecl{
		Address:   "0x7fd5e90e9078",
		Position:  "line:156:1, col:14",
		Position2: "col:14",
		Name:      "__stdinp",
		Type:      "FILE *",
		Type2:     "",
		IsExtern:  true,
		IsUsed:    false,
		IsCInit:   false,
		Children:  []interface{}{},
	},
	`0x7fd5e90ed630 <col:40, col:47> col:47 __size 'size_t':'unsigned long'`: &VarDecl{
		Address:   "0x7fd5e90ed630",
		Position:  "col:40, col:47",
		Position2: "col:47",
		Name:      "__size",
		Type:      "size_t",
		Type2:     "unsigned long",
		IsExtern:  false,
		IsUsed:    false,
		IsCInit:   false,
		Children:  []interface{}{},
	},
	`0x7fee35907a78 <col:4, col:8> col:8 used c 'int'`: &VarDecl{
		Address:   "0x7fee35907a78",
		Position:  "col:4, col:8",
		Position2: "col:8",
		Name:      "c",
		Type:      "int",
		Type2:     "",
		IsExtern:  false,
		IsUsed:    true,
		IsCInit:   false,
		Children:  []interface{}{},
	},
	`0x7fb0fd90ba30 <col:3, /usr/include/sys/_types.h:52:33> tests/assert/assert.c:13:9 used b 'int *' cinit`: &VarDecl{
		Address:   "0x7fb0fd90ba30",
		Position:  "col:3, /usr/include/sys/_types.h:52:33",
		Position2: "tests/assert/assert.c:13:9",
		Name:      "b",
		Type:      "int *",
		Type2:     "",
		IsExtern:  false,
		IsUsed:    true,
		IsCInit:   true,
		Children:  []interface{}{},
	},

	// WhileStmt
	`0x7fa1478273a0 <line:7:4, line:11:4>`: &WhileStmt{
		Address:  "0x7fa1478273a0",
		Position: "line:7:4, line:11:4",
		Children: []interface{}{},
	},
}

func TestNodes(t *testing.T) {
	for line, expected := range nodes {
		// Append the name of the struct onto the front. This would make
		// the complete line it would normally be parsing.
		actual := Parse(
			reflect.TypeOf(expected).Elem().Name() + " " + line)

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("\nexpected: %#v\n     got: %#v\n\n",
				expected, actual)
		}
	}
}
