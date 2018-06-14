package types_test

import (
	"encoding/json"
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
	{"int", "int32"},
	{"bool", "bool"},
	{"_Bool", "int8"},
	{"char", "byte"},
	{"char *[13]", "[]*byte"},
	{"__uint16_t", "uint16"},
	{"void *", "unsafe.Pointer"},
	{"unsigned short int", "uint16"},
	{"div_t", "noarch.DivT"},
	{"ldiv_t", "noarch.LdivT"},
	{"lldiv_t", "noarch.LldivT"},
	{"fpos_t", "int32"},
	{"int [2]", "[]int32"},
	{"int [2][3]", "[][]int32"},
	{"int [2][3][4]", "[][][]int32"},
	{"int [2][3][4][5]", "[][][][]int32"},
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
		{
			input:   "char *(*)( char *, ...)",
			fields:  []string{"char *", "..."},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)( char *, struct __va_list_tag *)",
			fields:  []string{"char *", "struct __va_list_tag *"},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)(const char *, ...)",
			fields:  []string{"const char *", "..."},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)(ImportCtx *)",
			fields:  []string{"ImportCtx *"},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)(int, char *, char *, ...)",
			fields:  []string{"int", "char *", "char *", "..."},
			returns: []string{"char *"},
		},
		{
			input:   "const char *(*)(int)",
			fields:  []string{"int"},
			returns: []string{"const char *"},
		},
		{
			input:   "const unsigned char *(*)(sqlite3_value *)",
			fields:  []string{"sqlite3_value *"},
			returns: []string{"const unsigned char *"},
		},
		{
			input:   "int (*)(const char *, sqlite3 **)",
			fields:  []string{"const char *", "sqlite3 **"},
			returns: []string{"int"},
		},
		{
			input: "int (*)(fts5_api *, const char *, void *, fts5_extension_function, void (*)(void *))",
			fields: []string{"fts5_api *",
				"const char *",
				"void *",
				"fts5_extension_function",
				"void (*)(void *)"},
			returns: []string{"int"},
		},
		{
			input: "int (*)(Fts5Context *, char *, int, void *, int (*)(void *, int, char *, int, int, int))",
			fields: []string{"Fts5Context *",
				"char *",
				"int",
				"void *",
				"int (*)(void *, int, char *, int, int, int)"},
			returns: []string{"int"},
		},
		{
			input: "int (*)(sqlite3 *, char *, int, int, void *, void (*)(sqlite3_context *, int, sqlite3_value **), void (*)(sqlite3_context *, int, sqlite3_value **), void (*)(sqlite3_context *))",
			fields: []string{
				"sqlite3 *",
				"char *",
				"int",
				"int",
				"void *",
				"void (*)(sqlite3_context *, int, sqlite3_value **)",
				"void (*)(sqlite3_context *, int, sqlite3_value **)",
				"void (*)(sqlite3_context *)",
			},
			returns: []string{"int"},
		}, /*
			{
				input: "int (*)(sqlite3_vtab *, int, const char *, void (**)(sqlite3_context *, int, sqlite3_value **), void **)",
				fields: []string{
					"sqlite3_vtab *",
					"int",
					"const char *",
					"void (**)(sqlite3_context *, int, sqlite3_value **)",
					"void **"},
				returns: []string{"int"},
			},*/
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
			if len(actualField) != len(tc.fields) {
				a, _ := json.Marshal(actualField)
				f, _ := json.Marshal(tc.fields)
				t.Errorf("Size of field is not same.\nActual  : %s\nExpected: %s\n",
					string(a),
					string(f))
				return
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
