package ast_test

import (
	"testing"
	"github.com/elliotchance/c2go/ast"
)

var nodes = map[string]interface{}{
	"0x7fce780f5018 </usr/include/sys/cdefs.h:313:68> always_inline":
	ast.AlwaysInlineAttr{
		"0x7fce780f5018",
		"/usr/include/sys/cdefs.h:313:68",
	},
	"0x7fe35b85d180 <col:63, col:69> 'char *' lvalue":
	ast.ArraySubscriptExpr{
		"0x7fe35b85d180",
		"col:63, col:69",
		"char *",
		"lvalue",
	},
}

func TestNodes(t *testing.T) {
	for line, expected := range nodes {
		var actual interface{}

		switch ty := expected.(type) {
		case ast.AlwaysInlineAttr:
			actual = ast.ParseAlwaysInlineAttr(line)
		case ast.ArraySubscriptExpr:
			actual = ast.ParseArraySubscriptExpr(line)
		default:
			t.Errorf("unknown %v", ty)
		}

		if expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	}
}
