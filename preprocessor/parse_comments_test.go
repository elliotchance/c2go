package preprocessor

import (
	"fmt"
	"testing"

	"github.com/elliotchance/c2go/program"
)

func TestParseComments(t *testing.T) {
	testCases := []struct {
		ent      entity
		code     []string
		comments []program.Comment
	}{
		{
			ent: entity{positionInSource: 10, include: "file.c"},
			code: []string{
				"NULL",
				"// comment1",
			},
			comments: []program.Comment{
				program.Comment{
					File:    "file.c",
					Line:    10,
					Comment: "// comment1",
				},
			},
		},
		{
			ent: entity{positionInSource: 10, include: "file.c"},
			code: []string{
				"NULL",
				"/* comment1 */",
			},
			comments: []program.Comment{
				program.Comment{
					File:    "file.c",
					Line:    10,
					Comment: "/* comment1 */",
				},
			},
		},
		{
			ent: entity{positionInSource: 10, include: "file.c"},
			code: []string{
				"NULL",
				"// comment1",
				"/* comment2 */",
			},
			comments: []program.Comment{
				program.Comment{
					File:    "file.c",
					Line:    10,
					Comment: "// comment1",
				},
				program.Comment{
					File:    "file.c",
					Line:    11,
					Comment: "/* comment2 */",
				},
			},
		},
		{
			ent: entity{positionInSource: 10, include: "file.c"},
			code: []string{
				"NULL",
				"/* comment2 */",
				"// comment1",
			},
			comments: []program.Comment{
				program.Comment{
					File:    "file.c",
					Line:    10,
					Comment: "/* comment2 */",
				},
				program.Comment{
					File:    "file.c",
					Line:    11,
					Comment: "// comment1",
				},
			},
		},
		{
			ent: entity{positionInSource: 10, include: "file.c"},
			code: []string{
				"NULL",
				"/* comment */ // comment1",
			},
			comments: []program.Comment{
				program.Comment{
					File:    "file.c",
					Line:    10,
					Comment: "/* comment */",
				},
				program.Comment{
					File:    "file.c",
					Line:    10,
					Comment: "// comment1",
				},
			},
		},
		{
			ent: entity{positionInSource: 10, include: "file.c"},
			code: []string{
				"NULL",
				"// comment1 /* comment */",
			},
			comments: []program.Comment{
				program.Comment{
					File:    "file.c",
					Line:    10,
					Comment: "// comment1 /* comment */",
				},
			},
		},
		{
			ent: entity{positionInSource: 10, include: "file.c"},
			code: []string{
				"NULL",
				"/* Text1",
				"Text2",
				"Text3 */",
			},
			comments: []program.Comment{
				program.Comment{
					File:    "file.c",
					Line:    10,
					Comment: "/* Text1\nText2\nText3 */",
				},
			},
		},
		{
			ent: entity{positionInSource: 10, include: "file.c"},
			code: []string{
				"NULL",
				"/* Text-1 */ // Text 0",
				"/* Text1",
				"Text2",
				"Text3 */ // Text 4",
				"// Text 5",
			},
			comments: []program.Comment{
				program.Comment{
					File:    "file.c",
					Line:    10,
					Comment: "/* Text-1 */",
				},
				program.Comment{
					File:    "file.c",
					Line:    10,
					Comment: "// Text 0",
				},
				program.Comment{
					File:    "file.c",
					Line:    11,
					Comment: "/* Text1\nText2\nText3 */",
				},
				program.Comment{
					File:    "file.c",
					Line:    13,
					Comment: "// Text 4",
				},
				program.Comment{
					File:    "file.c",
					Line:    14,
					Comment: "// Text 5",
				},
			},
		},
	}
	for i := range testCases {
		t.Run(fmt.Sprintf("Test:%d", i), func(t *testing.T) {
			var result []program.Comment
			for j := range testCases[i].code {
				testCases[i].ent.lines = append(testCases[i].ent.lines, &testCases[i].code[j])
			}
			testCases[i].ent.parseComments(&result)
			if len(result) != len(testCases[i].comments) {
				t.Fatalf("Size of comments is not same\nresult = '%d'\nexpect = '%d'",
					len(result),
					len(testCases[i].comments))
			}
			for j := range result {
				if result[j].File != testCases[i].comments[j].File {
					t.Fatalf("File is not same")
				}
				if result[j].Comment != testCases[i].comments[j].Comment {
					t.Fatalf("Comment is not same.\nresult = '%s'\nexpect = '%s'",
						result[j].Comment,
						testCases[i].comments[j].Comment)
				}
				if result[j].Line != testCases[i].comments[j].Line {
					t.Fatalf("Lines is not same\nresult = '%d'\nexpect = '%d'",
						result[j].Line,
						testCases[i].comments[j].Line)
				}
			}
		})
	}
}
