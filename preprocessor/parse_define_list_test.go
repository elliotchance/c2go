package preprocessor

import (
	"fmt"
	"testing"
)

func TestParseDefineList(t *testing.T) {
	testCases := []struct {
		inputLine string
		list      []string
	}{
		{
			inputLine: `#define A44 (p->u.s[24])
#define ArraySize(X) (int)(sizeof(X)/sizeof(X[0]))
#define BEGIN_TIMER beginTimer()
#define BIG_ENDIAN __BIG_ENDIAN
#define BUFSIZ _IO_BUFSIZ`,
			list: []string{
				"A44 (p->u.s[24])",
				"ArraySize(X) (int)(sizeof(X)/sizeof(X[0]))",
				"BEGIN_TIMER beginTimer()",
				"BIG_ENDIAN __BIG_ENDIAN",
				"BUFSIZ _IO_BUFSIZ",
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test:%d", i), func(t *testing.T) {
			actual, err := parseDefineList(tc.inputLine)
			if err != nil {
				t.Fatal(err)
			}
			if len(actual) != len(tc.list) {
				t.Fatalf("Cannot parse line : %s. Actual result : %#v. Expected: %#v", tc.inputLine, actual, tc.list)
			}
			for i := range actual {
				if actual[i] != tc.list[i] {
					t.Fatalf("Cannot parse '#define' in line : %s. Actual result : %#v. Expected: %#v", tc.inputLine, actual[i], tc.list[i])
				}
			}
		})
	}
}
