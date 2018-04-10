package ast

import (
	"fmt"
	"testing"
)

func TestFloatingLiteral(t *testing.T) {
	nodes := map[string]Node{
		`0x7febe106f5e8 <col:24> 'double' 1.230000e+00`: &FloatingLiteral{
			Addr:       0x7febe106f5e8,
			Pos:        NewPositionFromString("col:24"),
			Type:       "double",
			Value:      1.23,
			ChildNodes: []Node{},
		},
		`0x21c65b8 <col:41> 'double' 2.718282e+00`: &FloatingLiteral{
			Addr:       0x21c65b8,
			Pos:        NewPositionFromString("col:41"),
			Type:       "double",
			Value:      2.718282e+00,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}

func TestFloatingLiteralRepairFromSource(t *testing.T) {
	fl := &FloatingLiteral{
		Addr:       0x7febe106f5e8,
		Pos:        NewPositionFromString("col:12"),
		Type:       "double",
		Value:      1.23,
		ChildNodes: []Node{},
	}
	root := &CompoundStmt{
		Pos:        Position{File: "dummy.c", Line: 5},
		ChildNodes: []Node{fl},
	}
	FixPositions([]Node{root})
	type test struct {
		file     string
		expected float64
		err      error
	}
	tests := []test{
		{"# 2 \"x.c\"\n\n", 1.23, fmt.Errorf("could not find file %s", "dummy.c")},
		{"# 2 \"x.c\"\n\n# 1 \"dummy.c\"ff\nxxxxx\n\nyyyy", 1.23, fmt.Errorf("could not find %s:%d", "dummy.c", 5)},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = \nyyyy", 1.23, fmt.Errorf("cannot get exact value")},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = u\nyyyy", 1.23, fmt.Errorf("cannot parse float: strconv.ParseFloat: parsing \"\": invalid syntax from u")},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = 3.5zzz\nyyyy", 3.5, nil},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = .7\nyyyy", 0.7, nil},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = -.4e2\nyyyy", -40, nil},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = 0\nyyyy", 0, nil},
	}
	for _, test := range tests {
		prepareRepairFromSourceTest(t, test.file, func(ppFilePath string) {
			errors := RepairFloatingLiteralsFromSource(root, ppFilePath)
			if fl.Value != test.expected {
				t.Errorf("RepairFloatingLiteralsFromSource - expected: %f, got: %f", test.expected, fl.Value)
			}
			if test.err != nil && len(errors) == 0 || test.err == nil && len(errors) != 0 {
				t.Errorf("RepairFloatingLiteralsFromSource - error should match: expected: %v, got: %v", test.err, errors)
			} else if test.err != nil && errors[0].Err.Error() != test.err.Error() {
				t.Errorf("RepairFloatingLiteralsFromSource - error should match: expected: %s, got: %s", test.err.Error(), errors[0].Err.Error())
			}
		})
	}
}
