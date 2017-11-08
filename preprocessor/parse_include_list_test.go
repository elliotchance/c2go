package preprocessor

import (
	"fmt"
	"testing"
)

func TestParseIncludeList(t *testing.T) {
	testCases := []struct {
		inputLine string
		list      []string
	}{
		{
			inputLine: ` exit.o: exit.c tests.h `,
			list:      []string{"exit.c", "tests.h"},
		},
		{

			inputLine: ` exit.o: exit.c /usr/include/stdlib.h /usr/include/features.h \
  /usr/include/stdc-predef.h /usr/include/x86_64-linux-gnu/sys/cdefs.h \
  /usr/include/x86_64-linux-gnu/gnu/stubs-64.h \
  /usr/lib/llvm-3.8/bin/../lib/clang/3.8.0/include/stddef.h
  `,
			list: []string{"exit.c", "/usr/include/stdlib.h", "/usr/include/features.h",
				"/usr/include/stdc-predef.h", "/usr/include/x86_64-linux-gnu/sys/cdefs.h",
				"/usr/include/x86_64-linux-gnu/gnu/stubs-64.h",
				"/usr/lib/llvm-3.8/bin/../lib/clang/3.8.0/include/stddef.h",
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test:%d", i), func(t *testing.T) {
			actual, err := parseIncludeList(tc.inputLine)
			if err != nil {
				t.Fatal(err)
			}
			if len(actual) != len(tc.list) {
				t.Fatalf("Cannot parse line : %s. Actual result : %#v. Expected: %#v", tc.inputLine, actual, tc.list)
			}
			for i := range actual {
				if actual[i] != tc.list[i] {
					t.Fatalf("Cannot parse 'include' in line : %s. Actual result : %#v. Expected: %#v", tc.inputLine, actual[i], tc.list[i])
				}
			}
		})
	}
}
