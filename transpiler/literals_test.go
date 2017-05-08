package transpiler

import (
	"reflect"
	"testing"
	"unicode/utf8"

	"github.com/elliotchance/c2go/ast"
	goast "go/ast"
	"go/token"
)

var chartests = []struct {
	in  int    // Integer Character Code
	out string // Output Character Literal
}{
	// NUL byte
	{0, "'\\x00'"},

	// ASCII control characters
	{7, "'\\a'"},
	{8, "'\\b'"},
	{9, "'\\t'"},
	{10, "'\\n'"},
	{11, "'\\v'"},
	{12, "'\\f'"},
	{13, "'\\r'"},

	// printable ASCII
	{32, "' '"},
	{34, "'\"'"},
	{39, "'\\''"},
	{65, "'A'"},
	{92, "'\\\\'"},
	{191, "'¿'"},

	// printable unicode
	{948, "'δ'"},
	{0x03a9, "'Ω'"},
	{0x2020, "'†'"},

	// non-printable unicode
	{0xffff, "'\\uffff'"},
	{utf8.MaxRune, "'\\U0010ffff'"},
}

func TestCharacterLiterals(t *testing.T) {
	for _, tt := range chartests {
		expected := &goast.BasicLit{Kind: token.CHAR, Value: tt.out}
		actual := transpileCharacterLiteral(&ast.CharacterLiteral{Value: tt.in})
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("input: %v", tt.in)
			t.Errorf("  expected: %v", expected)
			t.Errorf("  actual:   %v", actual)
		}
	}
}
