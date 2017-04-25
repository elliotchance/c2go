package ast

import (
	"fmt"
	"reflect"
	"testing"
)

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
				t.Errorf("\nexpected: %#v\n     got: %#v\n\n",
					expected, actual)
			}
		})
	}
}
