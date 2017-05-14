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
