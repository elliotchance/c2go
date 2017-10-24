package preprocessor

import (
	"fmt"
	"testing"
)

func TestParseInclude(t *testing.T) {
	testCases := []struct {
		inputLine string
		include   string
	}{
		{
			inputLine: `# 1 "/usr/include/x86_64-linux-gnu/bits/sys_errlist.h" 1 3 4`,
			include:   "/usr/include/x86_64-linux-gnu/bits/sys_errlist.h",
		},
		{
			inputLine: `# 26 "/usr/include/x86_64-linux-gnu/bits/sys_errlist.h" 3 4`,
			include:   "/usr/include/x86_64-linux-gnu/bits/sys_errlist.h",
		},
		{
			inputLine: `# 854 "/usr/include/stdio.h" 2 3 4`,
			include:   "/usr/include/stdio.h",
		},
		{
			inputLine: `# 2 "f.c" 2`,
			include:   "f.c",
		},
		{
			inputLine: `# 30 "/usr/lib/llvm-3.8/bin/../lib/clang/3.8.0/include/stdarg.h" 3 4`,
			include:   "/usr/lib/llvm-3.8/bin/../lib/clang/3.8.0/include/stdarg.h",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test:%d", i), func(t *testing.T) {
			actual := parseInclude(tc.inputLine)
			if len(actual) == 0 {
				t.Fatal("Cannot parse, because result is empty")
			}
			if actual != tc.include {
				t.Fatalf("Cannot parse line: \"%s\". Result: \"%s\". Expected: \"%s\"", tc.inputLine, actual, tc.include)
			}
		})
	}
}
