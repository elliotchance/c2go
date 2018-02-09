package types_test

import (
	"fmt"
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
	{"void *", "interface{}"},
	{"unsigned short int", "uint16"},
	{"_Bool", "int"},
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

	for i, testCase := range resolveTestCases {
		t.Run(fmt.Sprintf("Test %d : %s", i, testCase.cType), func(t *testing.T) {
			goType, err := types.ResolveType(p, testCase.cType)
			if err != nil {
				t.Fatal(err)
			}

			if goType != testCase.goType {
				t.Errorf("Expected '%s' -> '%s', got '%s'",
					testCase.cType, testCase.goType, goType)
			}
		})
	}
}

func TestResolveFunction(t *testing.T) {
	var tcs = []struct {
		input   string
		fields  []string
		returns []string
	}{
		{
			input:   " void (*)(void)",
			fields:  []string{"void"},
			returns: []string{"void"},
		},
		{
			input:   " int (*)(sqlite3_file *)",
			fields:  []string{"sqlite3_file *"},
			returns: []string{"int"},
		},
		{
			input:   " int (*)(int)",
			fields:  []string{"int"},
			returns: []string{"int"},
		},
		{
			input:   " int (*)(void *) ",
			fields:  []string{"void *"},
			returns: []string{"int"},
		},
		{
			input:   " void (*)(sqlite3_context *, int, sqlite3_value **)",
			fields:  []string{"sqlite3_context *", "int", "sqlite3_value **"},
			returns: []string{"void"},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Test %d : %s", i, tc.input), func(t *testing.T) {
			actualField, actualReturn, err := types.ParseFunction(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if len(actualField) != len(tc.fields) {
				t.Error("Amount of fields is different")
			}
			for i := range actualField {
				if actualField[i] != tc.fields[i] {
					t.Errorf("Not correct field: %v\nExpected: %v", actualField, tc.fields)
				}
			}
			if len(actualReturn) != len(tc.returns) {
				t.Error("Amount of return elements are different")
			}
			for i := range actualReturn {
				if actualReturn[i] != tc.returns[i] {
					t.Errorf("Not correct returns: %v\nExpected: %v", actualReturn, tc.returns)
				}
			}
		})
	}

}
