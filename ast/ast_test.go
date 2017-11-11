package ast

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/elliotchance/c2go/util"
)

func formatMultiLine(o interface{}) string {
	s := fmt.Sprintf("%#v", o)
	s = strings.Replace(s, "{", "{\n", -1)
	s = strings.Replace(s, ", ", "\n", -1)

	return s
}

func runNodeTests(t *testing.T, tests map[string]Node) {
	i := 1
	for line, expected := range tests {
		testName := fmt.Sprintf("Example%d", i)
		i++

		t.Run(testName, func(t *testing.T) {
			// Append the name of the struct onto the front. This would make the
			// complete line it would normally be parsing.
			name := reflect.TypeOf(expected).Elem().Name()
			actual := Parse(name + " " + line)

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("%s", util.ShowDiff(formatMultiLine(expected),
					formatMultiLine(actual)))
			}
		})
	}
}

func TestPrint(t *testing.T) {
	cond := &ConditionalOperator{}
	cond.AddChild(&ImplicitCastExpr{})
	cond.AddChild(&ImplicitCastExpr{})
	s := Atos(cond)
	if len(s) == 0 {
		t.Fatalf("Cannot convert AST tree : %#v", cond)
	}
	lines := strings.Split(s, "\n")
	var amount int
	for _, l := range lines {
		if strings.Contains(l, "ImplicitCastExpr") {
			amount++
		}
	}
	if amount != 2 {
		t.Error("Not correct design of output")
	}
}

var lines = []string{
// c2go ast sqlite3.c | head -5000 | sed 's/^[ |`-]*//' | sed 's/<<<NULL>>>/NullStmt/g' | gawk 'length > 0 {print "`" $0 "`,"}'
}

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, line := range lines {
			Parse(line)
		}
	}
}
