// FIXME: Should be ast_test
package ast

import (
	"testing"
)

type resolveTestCase struct {
	cType  string
	goType string
}

var resolveTestCases = []resolveTestCase{
	{"int", "int"},
	{"char *[13]", "[]string"},
}

func TestResolve(t *testing.T) {
	for _, testCase := range resolveTestCases {
		goType := resolveType(testCase.cType)
		if goType != testCase.goType {
			t.Errorf("Expected '%s' -> '%s', got '%s'",
				testCase.cType, testCase.goType, goType)
		}
	}
}
