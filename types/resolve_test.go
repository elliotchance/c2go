package types_test

import (
	"testing"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
)

type resolveTestCase struct {
	cType  string
	goType string
}

var resolveTestCases = []resolveTestCase{
	{"int", "int"},
	{"char *[13]", "[][]byte"},
	{"__uint16_t", "uint16"},
	{"void *", "[]byte"},
	{"unsigned short int", "uint16"},
	{"_Bool", "bool"},
	{"struct RowSetEntry *", "[]RowSetEntry"},
	{"div_t", "noarch.DivT"},
	{"ldiv_t", "noarch.LdivT"},
	{"lldiv_t", "noarch.LldivT"},
	{"int [2]", "[]int"},
	{"int [2][3]", "[][]int"},
	{"int [2][3][4]", "[][][]int"},
	{"int [2][3][4][5]", "[][][][]int"},
}

func TestResolve(t *testing.T) {
	p := program.NewProgram()

	for _, testCase := range resolveTestCases {
		goType, err := types.ResolveType(p, testCase.cType)
		if err != nil {
			t.Error(err)
			continue
		}

		if goType != testCase.goType {
			t.Errorf("Expected '%s' -> '%s', got '%s'",
				testCase.cType, testCase.goType, goType)
		}
	}
}
