package types_test

import (
	"fmt"
	"testing"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
)

type sizeofTestCase struct {
	cType string
	size  int
	err   error
}

var sizeofTestCases = []sizeofTestCase{
	{"bool", 1, nil},
	{"_Bool", 1, nil},
	{"char", 1, nil},
	{"int", 4, nil},
	{"int [2]", 4 * 2, nil},
	{"int [2][3]", 4 * 2 * 3, nil},
	{"int [2][3][4]", 4 * 2 * 3 * 4, nil},
	{"int *[2]", 8 * 2, nil},
	{"int *[2][3]", 8 * 2 * 3, nil},
	{"int *[2][3][4]", 8 * 2 * 3 * 4, nil},
	{"int *", 8, nil},
	{"int **", 8, nil},
	{"int ***", 8, nil},
	{"char *const", 8, nil},
	{"char *const [3]", 24, nil},
	{"struct c [2]", 0, fmt.Errorf("Cannot determine sizeof : |struct c [2]|. err = error in sizeof baseSize")},
}

func TestSizeOf(t *testing.T) {
	p := program.NewProgram()

	for _, testCase := range sizeofTestCases {
		size, err := types.SizeOf(p, testCase.cType)
		if err != nil && (testCase.err == nil || (err.Error() != testCase.err.Error())) {
			t.Error(err)
			continue
		}

		if size != testCase.size {
			t.Errorf("Expected '%s' -> '%d', got '%d'",
				testCase.cType, testCase.size, size)
		}
	}
}
