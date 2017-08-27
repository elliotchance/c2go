package ast

import (
	"github.com/elliotchance/c2go/util"
	"regexp"
)

type Position struct {
	File      string // The relative or absolute file path.
	Line      int    // Start line
	LineEnd   int    // End line
	Column    int    // Start column
	ColumnEnd int    // End column

	// This is the original string that was converted. This is used for
	// debugging. We could derive this value from the other properties to save
	// on a bit of memory, but let worry about that later.
	StringValue string
}

func NewPositionFromString(s string) Position {
	if s == "<invalid sloc>" || s == "" {
		return Position{}
	}

	re := regexp.MustCompile(`^col:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Column:      util.Atoi(groups[1]),
		}
	}

	re = regexp.MustCompile(`^col:(\d+), col:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Column:      util.Atoi(groups[1]),
			ColumnEnd:   util.Atoi(groups[2]),
		}
	}

	re = regexp.MustCompile(`^line:(\d+), line:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Line:        util.Atoi(groups[1]),
			LineEnd:     util.Atoi(groups[2]),
		}
	}

	re = regexp.MustCompile(`^col:(\d+), line:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Column:      util.Atoi(groups[1]),
			Line:        util.Atoi(groups[2]),
		}
	}

	re = regexp.MustCompile(`^line:(\d+):(\d+), line:(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Line:        util.Atoi(groups[1]),
			Column:      util.Atoi(groups[2]),
			LineEnd:     util.Atoi(groups[3]),
			ColumnEnd:   util.Atoi(groups[4]),
		}
	}

	re = regexp.MustCompile(`^col:(\d+), line:(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Column:      util.Atoi(groups[1]),
			LineEnd:     util.Atoi(groups[2]),
			ColumnEnd:   util.Atoi(groups[3]),
		}
	}

	re = regexp.MustCompile(`^line:(\d+):(\d+), col:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Line:        util.Atoi(groups[1]),
			Column:      util.Atoi(groups[2]),
			ColumnEnd:   util.Atoi(groups[3]),
		}
	}

	re = regexp.MustCompile(`^line:(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		return Position{
			StringValue: s,
			Line:        util.Atoi(groups[1]),
			Column:      util.Atoi(groups[2]),
		}
	}

	// This must be below all of the others.
	re = regexp.MustCompile(`^([^:]+):(\d+):(\d+), col:(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		if groups[1] == "line" {
			panic(s)
		}
		return Position{
			StringValue: s,
			File:        groups[1],
			Line:        util.Atoi(groups[2]),
			Column:      util.Atoi(groups[3]),
			ColumnEnd:   util.Atoi(groups[4]),
		}
	}

	re = regexp.MustCompile(`^([^:]+):(\d+):(\d+), line:(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		if groups[1] == "line" {
			panic(s)
		}
		return Position{
			StringValue: s,
			File:        groups[1],
			Line:        util.Atoi(groups[2]),
			Column:      util.Atoi(groups[3]),
			LineEnd:     util.Atoi(groups[4]),
			ColumnEnd:   util.Atoi(groups[5]),
		}
	}

	re = regexp.MustCompile(`^([^:]+):(\d+):(\d+)$`)
	if groups := re.FindStringSubmatch(s); len(groups) > 0 {
		if groups[1] == "line" {
			panic(s)
		}
		return Position{
			StringValue: s,
			File:        groups[1],
			Line:        util.Atoi(groups[2]),
			Column:      util.Atoi(groups[3]),
		}
	}

	panic("unable to understand position '" + s + "'")
}
