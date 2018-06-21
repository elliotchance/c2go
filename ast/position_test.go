package ast

import (
	"testing"
)

func TestNewPositionFromString(t *testing.T) {
	tests := map[string]Position{
		`col:30`: {
			File:      "",
			Line:      0,
			Column:    30,
			LineEnd:   0,
			ColumnEnd: 0,
		},
		`col:47, col:57`: {
			File:      "",
			Line:      0,
			Column:    47,
			LineEnd:   0,
			ColumnEnd: 57,
		},
		`/usr/include/sys/cdefs.h:313:68`: {
			File:      "/usr/include/sys/cdefs.h",
			Line:      313,
			Column:    68,
			LineEnd:   0,
			ColumnEnd: 0,
		},
		`C:\usr\include\sys\cdefs.h:313:68`: {
			File:      `C:\usr\include\sys\cdefs.h`,
			Line:      313,
			Column:    68,
			LineEnd:   0,
			ColumnEnd: 0,
		},
		`/usr/include/AvailabilityInternal.h:21697:88, col:124`: {
			File:      "/usr/include/AvailabilityInternal.h",
			Line:      21697,
			Column:    88,
			LineEnd:   0,
			ColumnEnd: 124,
		},
		`C:\usr\include\AvailabilityInternal.h:21697:88, col:124`: {
			File:      `C:\usr\include\AvailabilityInternal.h`,
			Line:      21697,
			Column:    88,
			LineEnd:   0,
			ColumnEnd: 124,
		},
		`line:275:50, col:99`: {
			File:      "",
			Line:      275,
			Column:    50,
			LineEnd:   0,
			ColumnEnd: 99,
		},
		`line:11:5, line:12:21`: {
			File:      "",
			Line:      11,
			Column:    5,
			LineEnd:   12,
			ColumnEnd: 21,
		},
		`col:54, line:358:1`: {
			File:      "",
			Line:      0,
			Column:    54,
			LineEnd:   358,
			ColumnEnd: 1,
		},
		`/usr/include/secure/_stdio.h:42:1, line:43:32`: {
			File:      "/usr/include/secure/_stdio.h",
			Line:      42,
			Column:    1,
			LineEnd:   43,
			ColumnEnd: 32,
		},
		`C:\usr\include\secure\_stdio.h:42:1, line:43:32`: {
			File:      `C:\usr\include\secure\_stdio.h`,
			Line:      42,
			Column:    1,
			LineEnd:   43,
			ColumnEnd: 32,
		},
		`line:244:5`: {
			File:      "",
			Line:      244,
			Column:    5,
			LineEnd:   0,
			ColumnEnd: 0,
		},
		`<invalid sloc>`: {
			File:      "",
			Line:      0,
			Column:    0,
			LineEnd:   0,
			ColumnEnd: 0,
		},
		`/usr/include/sys/stdio.h:39:1, /usr/include/AvailabilityInternal.h:21697:126`: {
			File:      "/usr/include/sys/stdio.h",
			Line:      39,
			Column:    1,
			LineEnd:   0,
			ColumnEnd: 0,
		},
		`C:\usr\include\sys\stdio.h:39:1, C:\usr\include\AvailabilityInternal.h:21697:126`: {
			File:      `C:\usr\include\sys\stdio.h`,
			Line:      39,
			Column:    1,
			LineEnd:   0,
			ColumnEnd: 0,
		},
		`C:\Program Files (x86)\Microsoft Visual Studio 11.0\VC\include\math.h:39:1, C:\Program Files (x86)\Microsoft Visual Studio 11.0\VC\include\AvailabilityInternal.h:21697:126`: {
			File:      `C:\Program Files (x86)\Microsoft Visual Studio 11.0\VC\include\math.h`,
			Line:      39,
			Column:    1,
			LineEnd:   0,
			ColumnEnd: 0,
		},
		`col:1, /usr/include/sys/cdefs.h:351:63`: {
			File:      "",
			Line:      0,
			Column:    1,
			LineEnd:   0,
			ColumnEnd: 0,
		},
		`col:1, C:\usr\include\sys\cdefs.h:351:63`: {
			File:      "",
			Line:      0,
			Column:    1,
			LineEnd:   0,
			ColumnEnd: 0,
		},
	}

	for testName, expectedPos := range tests {
		t.Run(testName, func(t *testing.T) {
			pos := NewPositionFromString(testName)
			if pos.File != expectedPos.File {
				t.Errorf("TestNewPositionFromString: File: %#v != %#v", pos.File, expectedPos.File)
			}
			if pos.Line != expectedPos.Line {
				t.Errorf("TestNewPositionFromString: Line: %#v != %#v", pos.Line, expectedPos.Line)
			}
			if pos.Column != expectedPos.Column {
				t.Errorf("TestNewPositionFromString: Column: %#v != %#v", pos.Column, expectedPos.Column)
			}
		})
	}
}
