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
