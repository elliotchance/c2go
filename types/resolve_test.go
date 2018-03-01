package types_test

import (
	"encoding/json"
	"fmt"
	"strings"
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
	{"div_t", "noarch.DivT"},
	{"ldiv_t", "noarch.LdivT"},
	{"lldiv_t", "noarch.LldivT"},
	{"int [2]", "[]int"},
	{"int [2][3]", "[][]int"},
	{"int [2][3][4]", "[][][]int"},
	{"int [2][3][4][5]", "[][][][]int"},
	{"int (*[2])(int, int)", "[2]func(int,int)(int)"},
	{"int (*(*(*)))(int, int)", "[][]func(int,int)(int)"},
}

func TestResolve(t *testing.T) {
	p := program.NewProgram()

	for i, testCase := range resolveTestCases {
		t.Run(fmt.Sprintf("Test %d : %s", i, testCase.cType), func(t *testing.T) {
			goType, err := types.ResolveType(p, testCase.cType)
			if err != nil {
				t.Fatal(err)
			}

			goType = strings.Replace(goType, " ", "", -1)
			testCase.goType = strings.Replace(testCase.goType, " ", "", -1)

			if goType != testCase.goType {
				t.Errorf("Expected '%s' -> '%s', got '%s'",
					testCase.cType, testCase.goType, goType)
			}
		})
	}
}

func TestResolveFunction(t *testing.T) {
	var tcs = []struct {
		input string

		prefix  string
		fields  []string
		returns []string
	}{
		{
			input:   "__ssize_t (void *, char *, size_t)",
			prefix:  "",
			fields:  []string{"void *", "char *", "size_t"},
			returns: []string{"__ssize_t"},
		},
		{
			input:  "int (*)(sqlite3_vtab *, int, const char *, void (**)(sqlite3_context *, int, sqlite3_value **), void **)",
			prefix: "",
			fields: []string{
				"sqlite3_vtab *",
				"int",
				"const char *",
				"void (**)(sqlite3_context *, int, sqlite3_value **)",
				"void **"},
			returns: []string{"int"},
		},
		{
			input:   "void ( *(*)(int *, void *, char *))(void)",
			prefix:  "*",
			fields:  []string{"void"},
			returns: []string{"void (int *, void *, char *)"},
		},
		{
			input:   " void (*)(void)",
			prefix:  "",
			fields:  []string{"void"},
			returns: []string{"void"},
		},
		{
			input:   " int (*)(sqlite3_file *)",
			prefix:  "",
			fields:  []string{"sqlite3_file *"},
			returns: []string{"int"},
		},
		{
			input:   " int (*)(int)",
			prefix:  "",
			fields:  []string{"int"},
			returns: []string{"int"},
		},
		{
			input:   " int (*)(void *) ",
			prefix:  "",
			fields:  []string{"void *"},
			returns: []string{"int"},
		},
		{
			input:   " void (*)(sqlite3_context *, int, sqlite3_value **)",
			prefix:  "",
			fields:  []string{"sqlite3_context *", "int", "sqlite3_value **"},
			returns: []string{"void"},
		},
		{
			input:   "char *(*)( char *, ...)",
			prefix:  "",
			fields:  []string{"char *", "..."},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)( char *, struct __va_list_tag *)",
			prefix:  "",
			fields:  []string{"char *", "struct __va_list_tag *"},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)(const char *, ...)",
			prefix:  "",
			fields:  []string{"const char *", "..."},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)(ImportCtx *)",
			prefix:  "",
			fields:  []string{"ImportCtx *"},
			returns: []string{"char *"},
		},
		{
			input:   "char *(*)(int, char *, char *, ...)",
			prefix:  "",
			fields:  []string{"int", "char *", "char *", "..."},
			returns: []string{"char *"},
		},
		{
			input:   "const char *(*)(int)",
			prefix:  "",
			fields:  []string{"int"},
			returns: []string{"const char *"},
		},
		{
			input:   "const unsigned char *(*)(sqlite3_value *)",
			prefix:  "",
			fields:  []string{"sqlite3_value *"},
			returns: []string{"const unsigned char *"},
		},
		{
			input:   "int (*)(const char *, sqlite3 **)",
			prefix:  "",
			fields:  []string{"const char *", "sqlite3 **"},
			returns: []string{"int"},
		},
		{
			input:  "int (*)(fts5_api *, const char *, void *, fts5_extension_function, void (*)(void *))",
			prefix: "",
			fields: []string{"fts5_api *",
				"const char *",
				"void *",
				"fts5_extension_function",
				"void (*)(void *)"},
			returns: []string{"int"},
		},
		{
			input:  "int (*)(Fts5Context *, char *, int, void *, int (*)(void *, int, char *, int, int, int))",
			prefix: "",
			fields: []string{"Fts5Context *",
				"char *",
				"int",
				"void *",
				"int (*)(void *, int, char *, int, int, int)"},
			returns: []string{"int"},
		},
		{
			input:  "int (*)(sqlite3 *, char *, int, int, void *, void (*)(sqlite3_context *, int, sqlite3_value **), void (*)(sqlite3_context *, int, sqlite3_value **), void (*)(sqlite3_context *))",
			prefix: "",
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
		},
		{
			input:  "int (*)(sqlite3_vtab *, int, const char *, void (**)(sqlite3_context *, int, sqlite3_value **), void **)",
			prefix: "",
			fields: []string{
				"sqlite3_vtab *",
				"int",
				"const char *",
				"void (**)(sqlite3_context *, int, sqlite3_value **)",
				"void **"},
			returns: []string{"int"},
		},
		{
			input:   "void (*(int *, void *, const char *))(void)",
			prefix:  "",
			fields:  []string{"void"},
			returns: []string{"void (int *, void *, const char *)"},
		},
		{
			input:   "long (int, int)",
			prefix:  "",
			fields:  []string{"int", "int"},
			returns: []string{"long"},
		},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Test %d : %s", i, tc.input), func(t *testing.T) {
			actualPrefix, actualField, actualReturn, err :=
				types.ParseFunction(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if actualPrefix != tc.prefix {
				t.Errorf("Prefix is not same.\nActual: %s\nExpected: %s\n",
					actualPrefix, tc.prefix)
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
				actualField[i] = strings.Replace(actualField[i], " ", "", -1)
				tc.fields[i] = strings.Replace(tc.fields[i], " ", "", -1)
				if actualField[i] != tc.fields[i] {
					t.Errorf("Not correct field: %v\nExpected: %v", actualField, tc.fields)
				}
			}
			if len(actualReturn) != len(tc.returns) {
				a, _ := json.Marshal(actualReturn)
				f, _ := json.Marshal(tc.returns)
				t.Errorf("Size of return field is not same.\nActual  : %s\nExpected: %s\n",
					string(a),
					string(f))
				return
			}
			if len(actualReturn) != len(tc.returns) {
				t.Errorf("Amount of return elements are different\nActual  : %v\nExpected: %v\n",
					actualReturn, tc.returns)
			}
			for i := range actualReturn {
				actualReturn[i] = strings.Replace(actualReturn[i], " ", "", -1)
				tc.returns[i] = strings.Replace(tc.returns[i], " ", "", -1)
				if actualReturn[i] != tc.returns[i] {
					t.Errorf("Not correct returns: %v\nExpected: %v", actualReturn, tc.returns)
				}
			}
		})
	}

}
