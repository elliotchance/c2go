package preprocessor

import (
	"fmt"
	"testing"
)

func TestParseIncludePreproccessorLine(t *testing.T) {
	testCases := []struct {
		inputLine string
		out       entity
	}{
		{
			inputLine: `# 1 "/usr/include/x86_64-linux-gnu/bits/sys_errlist.h" 1 3 4`,
			out: entity{
				include:          "/usr/include/x86_64-linux-gnu/bits/sys_errlist.h",
				positionInSource: 1,
			},
		},
		{
			inputLine: `# 26 "/usr/include/x86_64-linux-gnu/bits/sys_errlist.h" 3 4`,
			out: entity{
				include:          "/usr/include/x86_64-linux-gnu/bits/sys_errlist.h",
				positionInSource: 26,
			},
		},
		{
			inputLine: `# 854 "/usr/include/stdio.h" 2 3 4`,
			out: entity{
				include:          "/usr/include/stdio.h",
				positionInSource: 854,
			},
		},
		{
			inputLine: `# 2 "f.c" 2`,
			out: entity{
				include:          "f.c",
				positionInSource: 2,
			},
		},
		{
			inputLine: `# 2 "f.c"`,
			out: entity{
				include:          "f.c",
				positionInSource: 2,
			},
		},
		{
			inputLine: `# 30 "/usr/lib/llvm-3.8/bin/../lib/clang/3.8.0/include/stdarg.h" 3 4`,
			out: entity{
				include:          "/usr/lib/llvm-3.8/bin/../lib/clang/3.8.0/include/stdarg.h",
				positionInSource: 30,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test:%d", i), func(t *testing.T) {
			actual, err := parseIncludePreprocessorLine(tc.inputLine)
			if err != nil {
				t.Fatal(err)
			}
			if len(actual.include) == 0 {
				t.Fatal("Cannot parse, because result is empty")
			}
			if actual.include != tc.out.include {
				t.Fatalf("Cannot parse line: \"%s\". Result: \"%s\". Expected: \"%s\"", tc.inputLine, actual.include, tc.out.include)
			}
			if actual.positionInSource != tc.out.positionInSource {
				t.Fatalf("Cannot parse source position in line: \"%s\". Result: \"%d\". Expected: \"%d\"", tc.inputLine, actual.positionInSource, tc.out.positionInSource)
			}
		})
	}
}

func TestParseIncludePreproccessorLineFail1(t *testing.T) {
	inputLine := `# A "/usr/include/stdio.h" 2 3 4`
	_, err := parseIncludePreprocessorLine(inputLine)
	if err == nil {
		t.Fatal("Cannot found error in positionInSource")
	}
}

func TestParseIncludePreproccessorLineFail2(t *testing.T) {
	inputLine := ` # 850 "/usr/include/stdio.h" 2 3 4`
	_, err := parseIncludePreprocessorLine(inputLine)
	if err == nil {
		t.Fatal("Cannot give error if first symbol is not #")
	}
}

func TestParseIncludePreproccessorLineFail3(t *testing.T) {
	inputLine := `# 850`
	_, err := parseIncludePreprocessorLine(inputLine)
	if err == nil {
		t.Fatal("Cannot give error if line hanen't include string")
	}
}

func TestParseIncludePreproccessorLineFail4(t *testing.T) {
	inputLine := `# 850 "/usr/include`
	_, err := parseIncludePreprocessorLine(inputLine)
	if err == nil {
		t.Fatal("Cannot give error if wrong format of include line")
	}
}

func TestParseIncludePreproccessorLineFail5(t *testing.T) {
	inputLine := `# 850`
	_, err := parseIncludePreprocessorLine(inputLine)
	if err == nil {
		t.Fatal("Cannot give error if haven`t include line")
	}
}
