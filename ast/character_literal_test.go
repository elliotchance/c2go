package ast

import (
	"fmt"
	"github.com/elliotchance/c2go/cc"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestCharacterLiteral(t *testing.T) {
	nodes := map[string]Node{
		`0x7f980b858308 <col:62> 'int' 10`: &CharacterLiteral{
			Addr:       0x7f980b858308,
			Pos:        NewPositionFromString("col:62"),
			Type:       "int",
			Value:      10,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}

func TestCharacterLiteralDecodeHex(t *testing.T) {
	hexCodes := map[rune]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'a': 10,
		'b': 11,
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
		'A': 10,
		'B': 11,
		'C': 12,
		'D': 13,
		'E': 14,
		'F': 15,
		'z': -1, // illegal hex characters return -1
		'G': -1,
	}
	for r, expected := range hexCodes {
		i := decodeHex(r)

		if i != expected {
			t.Errorf("hexDecode should match: expected: %d, got: %d", expected, i)
		}
	}
}

func TestCharacterLiteralDecodeHexQuad(t *testing.T) {
	runes := map[string][]int{
		"zzzzzz": {0, 0},
		"1zzzzz": {0x1, 1},
		"1Azzzz": {0x1a, 2},
		"1A3zzz": {0x1a3, 3},
		"1A3bzz": {0x1a3b, 4},
		"1A3b5z": {0x1a3b, 4},
		"1A3b":   {0x1a3b, 4},
		"1A3":    {0x1a3, 3},
		"1A":     {0x1a, 2},
		"1":      {0x1, 1},
		"":       {0, 0},
	}
	for r, expected := range runes {
		i, n := decodeHexQuad([]rune(r))

		if int(i) != expected[0] {
			t.Errorf("hexDecodeQuad('%s') should match: expected: %x, got: %x", r, expected[0], i)
		}
		if n != expected[1] {
			t.Errorf("hexDecodeQuad('%s') length should match: expected: %d, got: %d", r, expected[1], n)
		}
	}
}

func TestCharacterLiteralDecodeUCN(t *testing.T) {
	runes := map[string][]int{
		"\\uzzzzzz":     {0, 2},
		"\\u1zzzzz":     {0x1, 3},
		"\\u1Azzzz":     {0x1a, 4},
		"\\u1A3zzz":     {0x1a3, 5},
		"\\u1A3bzz":     {0x1a3b, 6},
		"\\u1A3b5z":     {0x1a3b, 6},
		"\\u1A3b":       {0x1a3b, 6},
		"\\u1A3":        {0x1a3, 5},
		"\\u1A":         {0x1a, 4},
		"\\u1":          {0x1, 3},
		"\\u":           {0, 2},
		"\\Uzzzzzzzzzz": {0, 2},
		"\\U1zzzzzzzzz": {0x1, 3},
		"\\U1Azzzzzzzz": {0x1a, 4},
		"\\U1A3zzzzzzz": {0x1a3, 5},
		"\\U1A3bzzzzzz": {0x1a3b, 6},
		"\\U1A3b5zzzzz": {0x1a3b5, 7},
		"\\U1A3b52zzzz": {0x1a3b52, 8},
		"\\U1A3b52fzzz": {0x1a3b52f, 9},
		"\\U1A3b52f6zz": {0x1a3b52f6, 10},
		"\\U1A3b52f69z": {0x1a3b52f6, 10},
		"\\U1A3b52f6":   {0x1a3b52f6, 10},
		"\\U1A3b52f":    {0x1a3b52f, 9},
		"\\U1A3b52":     {0x1a3b52, 8},
		"\\U1A3b5":      {0x1a3b5, 7},
		"\\U1A3b":       {0x1a3b, 6},
		"\\U1A3":        {0x1a3, 5},
		"\\U1A":         {0x1a, 4},
		"\\U1":          {0x1, 3},
		"\\U":           {0, 2},
	}
	for r, expected := range runes {
		i, n := decodeUCN([]rune(r))

		if int(i) != expected[0] {
			t.Errorf("decodeUCN('%s') should match: expected: %x, got: %x", r, expected[0], int(i))
		}
		if n != expected[1] {
			t.Errorf("decodeUCN('%s') length should match: expected: %d, got: %d", r, expected[1], n)
		}
	}
}

func TestCharacterLiteralDecodeEscapeSequence(t *testing.T) {
	type test struct {
		result int
		length int
		err    error
	}
	runes := map[string]test{
		"\\'zz":         {int('\''), 2, nil},
		"\\\"zz":        {int('"'), 2, nil},
		"\\?zz":         {int('?'), 2, nil},
		"\\\\zz":        {int('\\'), 2, nil},
		"\\azz":         {int('\a'), 2, nil},
		"\\bzz":         {int('\b'), 2, nil},
		"\\fzz":         {int('\f'), 2, nil},
		"\\nzz":         {int('\n'), 2, nil},
		"\\rzz":         {int('\r'), 2, nil},
		"\\tzz":         {int('\t'), 2, nil},
		"\\vzz":         {int('\v'), 2, nil},
		"\\'":           {int('\''), 2, nil},
		"\\\"":          {int('"'), 2, nil},
		"\\?":           {int('?'), 2, nil},
		"\\\\":          {int('\\'), 2, nil},
		"\\a":           {int('\a'), 2, nil},
		"\\b":           {int('\b'), 2, nil},
		"\\f":           {int('\f'), 2, nil},
		"\\n":           {int('\n'), 2, nil},
		"\\r":           {int('\r'), 2, nil},
		"\\t":           {int('\t'), 2, nil},
		"\\v":           {int('\v'), 2, nil},
		"\\xzzzzzz":     {0, 2, nil},
		"\\x3zzzzz":     {0x3, 3, nil},
		"\\x3bzzzz":     {0x3b, 4, nil},
		"\\x3b1zzz":     {0x3b, 4, nil},
		"\\x3b":         {0x3b, 4, nil},
		"\\x3":          {0x3, 3, nil},
		"\\x":           {0, 2, nil},
		"\\zzzzzz":      {0, 0, fmt.Errorf("illegal character '%s'", "z")},
		"\\dzzzzz":      {0, 0, fmt.Errorf("illegal character '%s'", "d")},
		"\\9zzzzz":      {0, 0, fmt.Errorf("illegal character '%s'", "9")},
		"\\z":           {0, 0, fmt.Errorf("illegal character '%s'", "z")},
		"\\d":           {0, 0, fmt.Errorf("illegal character '%s'", "d")},
		"\\9":           {0, 0, fmt.Errorf("illegal character '%s'", "9")},
		"\\0zzzzzz":     {0, 2, nil},
		"\\1zzzzzz":     {1, 2, nil},
		"\\2zzzzzz":     {2, 2, nil},
		"\\3zzzzzz":     {3, 2, nil},
		"\\4zzzzzz":     {4, 2, nil},
		"\\5zzzzzz":     {5, 2, nil},
		"\\6zzzzzz":     {6, 2, nil},
		"\\7zzzzzz":     {7, 2, nil},
		"\\0":           {0, 2, nil},
		"\\1":           {1, 2, nil},
		"\\2":           {2, 2, nil},
		"\\3":           {3, 2, nil},
		"\\4":           {4, 2, nil},
		"\\5":           {5, 2, nil},
		"\\6":           {6, 2, nil},
		"\\7":           {7, 2, nil},
		"\\34zzzzzz":    {034, 3, nil},
		"\\347zzzzz":    {0347, 4, nil},
		"\\3472zzzz":    {0347, 4, nil},
		"\\347":         {0347, 4, nil},
		"\\34":          {034, 3, nil},
		"\\uzzzzzz":     {0, 2, nil},
		"\\u1zzzzz":     {0x1, 3, nil},
		"\\u1Azzzz":     {0x1a, 4, nil},
		"\\u1A3zzz":     {0x1a3, 5, nil},
		"\\u1A3bzz":     {0x1a3b, 6, nil},
		"\\u1A3b5z":     {0x1a3b, 6, nil},
		"\\u1A3b":       {0x1a3b, 6, nil},
		"\\u1A3":        {0x1a3, 5, nil},
		"\\u1A":         {0x1a, 4, nil},
		"\\u1":          {0x1, 3, nil},
		"\\u":           {0, 2, nil},
		"\\Uzzzzzzzzzz": {0, 2, nil},
		"\\U1zzzzzzzzz": {0x1, 3, nil},
		"\\U1Azzzzzzzz": {0x1a, 4, nil},
		"\\U1A3zzzzzzz": {0x1a3, 5, nil},
		"\\U1A3bzzzzzz": {0x1a3b, 6, nil},
		"\\U1A3b5zzzzz": {0x1a3b5, 7, nil},
		"\\U1A3b52zzzz": {0x1a3b52, 8, nil},
		"\\U1A3b52fzzz": {0x1a3b52f, 9, nil},
		"\\U1A3b52f6zz": {0x1a3b52f6, 10, nil},
		"\\U1A3b52f69z": {0x1a3b52f6, 10, nil},
		"\\U1A3b52f6":   {0x1a3b52f6, 10, nil},
		"\\U1A3b52f":    {0x1a3b52f, 9, nil},
		"\\U1A3b52":     {0x1a3b52, 8, nil},
		"\\U1A3b5":      {0x1a3b5, 7, nil},
		"\\U1A3b":       {0x1a3b, 6, nil},
		"\\U1A3":        {0x1a3, 5, nil},
		"\\U1A":         {0x1a, 4, nil},
		"\\U1":          {0x1, 3, nil},
		"\\U":           {0, 2, nil},
	}
	for r, expected := range runes {
		i, n, err := decodeEscapeSequence([]rune(r))

		if int(i) != expected.result {
			t.Errorf("decodeEscapeSequence('%s') should match: expected: %x, got: %x", r, expected.result, int(i))
		}
		if n != expected.length {
			t.Errorf("decodeEscapeSequence('%s') length should match: expected: %d, got: %d", r, expected.length, n)
		}
		if err != nil && expected.err == nil || err == nil && expected.err != nil {
			t.Errorf("decodeEscapeSequence('%s') error should match: expected: %v, got: %v", r, expected.err, err)
		} else if err != nil && err.Error() != expected.err.Error() {
			t.Errorf("decodeEscapeSequence('%s') error should match: expected: %s, got: %s", r, expected.err.Error(), err.Error())
		}
	}
}

func TestCharacterLiteralFromSource(t *testing.T) {
	type test struct {
		result int
		err    error
	}
	prefixes := map[string]error{
		"'":   nil,
		"u'":  nil,
		"u8'": nil,
		"U'":  nil,
		"L'":  nil,
		"z":   fmt.Errorf("illegal character 'z' at index 0"),
		"uz":  fmt.Errorf("illegal character 'z' at index 1"),
		"Uz":  fmt.Errorf("illegal character 'z' at index 1"),
		"Lz":  fmt.Errorf("illegal character 'z' at index 1"),
		"u8z": fmt.Errorf("illegal character 'z' at index 2"),
	}
	suffixes := map[string]test{
		"":     {0, fmt.Errorf("unexpected end of character literal")},
		"'":    {0, fmt.Errorf("empty character literal")},
		"a'":   {int('a'), nil},
		"a":    {0, fmt.Errorf("unexpected end of character literal")},
		"\\":   {0, fmt.Errorf("unexpected end of character literal")},
		"\\'":  {0, fmt.Errorf("unexpected end of character literal")},
		"\\''": {int('\''), nil},
		"ab":   {0, fmt.Errorf("does not support multi-character literals")},
		"\\nb": {0, fmt.Errorf("does not support multi-character literals")},
	}
	tests := map[string]test{
		"":   {0, fmt.Errorf("character literal to short")},
		"u":  {0, fmt.Errorf("character literal to short")},
		"u8": {0, fmt.Errorf("character literal to short")},
		"U":  {0, fmt.Errorf("character literal to short")},
		"L":  {0, fmt.Errorf("character literal to short")},
	}
	for i, x := range prefixes {
		for j, y := range suffixes {
			var t = y
			str := i + j
			if x != nil {
				t.result = 0
				t.err = x
			}
			tests[str] = t
		}
	}
	for s, expected := range tests {
		i, err := parseCharacterLiteralFromSource(s)

		if int(i) != expected.result {
			t.Errorf("parseCharacterLiteralFromSource('%s') should match: expected: %x, got: %x", s, expected.result, int(i))
		}
		if err != nil && expected.err == nil || err == nil && expected.err != nil {
			t.Errorf("parseCharacterLiteralFromSource('%s') error should match: expected: %v, got: %v", s, expected.err, err)
		} else if err != nil && err.Error() != expected.err.Error() {
			t.Errorf("parseCharacterLiteralFromSource('%s') error should match: expected: %s, got: %s", s, expected.err.Error(), err.Error())
		}
	}
}

func TestCharacterLiteralRepairFromSource(t *testing.T) {
	cl := &CharacterLiteral{
		Addr:       0x7f980b858308,
		Pos:        NewPositionFromString("col:12"),
		Type:       "int",
		Value:      10,
		ChildNodes: []Node{},
	}
	root := &CompoundStmt{
		Pos:        Position{File: "dummy.c", Line: 5},
		ChildNodes: []Node{cl},
	}
	FixPositions([]Node{root})
	type test struct {
		file     string
		expected int
		err      error
	}
	tests := []test{
		{"# 2 \"x.c\"\n\n", 10, fmt.Errorf("could not find file %s", "dummy.c")},
		{"# 2 \"x.c\"\n\n# 1 \"dummy.c\"ff\nxxxxx\n\nyyyy", 10, fmt.Errorf("could not find %s:%d", "dummy.c", 5)},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = \nyyyy", 10, fmt.Errorf("cannot get exact value")},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = u\nyyyy", 10, fmt.Errorf("cannot parse character literal: character literal to short from u")},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = '\nyyyy", 10, fmt.Errorf("cannot parse character literal: unexpected end of character literal from '")},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = ''zzz\nyyyy", 10, fmt.Errorf("cannot parse character literal: empty character literal from ''zzz")},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = 'A'zzz\nyyyy", int('A'), nil},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = 'B'\nyyyy", int('B'), nil},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = '\\n'\nyyyy", int('\n'), nil},
	}
	for _, test := range tests {
		prepareRepairFromSourceTest(t, test.file, func(ppFilePath string) {
			errors := RepairCharacterLiteralsFromSource(root, ppFilePath)
			if cl.Value != test.expected {
				t.Errorf("RepairCharacterLiteralsFromSource - expected: %x, got: %x", test.expected, cl.Value)
			}
			if test.err != nil && len(errors) == 0 || test.err == nil && len(errors) != 0 {
				t.Errorf("RepairCharacterLiteralsFromSource - error should match: expected: %v, got: %v", test.err, errors)
			} else if test.err != nil && errors[0].Err.Error() != test.err.Error() {
				t.Errorf("RepairCharacterLiteralsFromSource - error should match: expected: %s, got: %s", test.err.Error(), errors[0].Err.Error())
			}
		})
	}
}

func prepareRepairFromSourceTest(t *testing.T, fileContent string, test func(filePath string)) {
	cc.ResetCache()
	dir, err := ioutil.TempDir("", "c2go")
	if err != nil {
		t.Fatal(fmt.Errorf("Cannot create temp folder: %v", err))
	}
	defer os.RemoveAll(dir) // clean up

	ppFilePath := path.Join(dir, "pp.c")
	err = ioutil.WriteFile(ppFilePath, []byte(fileContent), 0644)
	if err != nil {
		t.Fatal(fmt.Errorf("writing to %s failed: %v", ppFilePath, err))
	}

	test(ppFilePath)
}

func TestCharacterLiteralRepairFromSourceMultiline(t *testing.T) {
	cl := &CharacterLiteral{
		Addr:       0x7f980b858308,
		Pos:        NewPositionFromString("col:12"),
		Type:       "int",
		Value:      10,
		ChildNodes: []Node{},
	}
	cl.Pos.LineEnd = 6
	root := &CompoundStmt{
		Pos:        Position{File: "dummy.c", Line: 5},
		ChildNodes: []Node{cl},
	}
	FixPositions([]Node{root})
	type test struct {
		file     string
		expected int
		err      error
	}
	tests := []test{
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = '\nxxxxxxxxx, };\nyyyy", 10, fmt.Errorf("cannot parse character literal: illegal character '}' at index 0 from };")},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = 'A'zzz\nyyyy", int('A'), nil},
		{"# 2 \"x.c\"\n\n# 4 \"dummy.c\"ff\nxxxxx\nvar xyst = {\nxxxxxxxxx, 'B'};\nyyyy", int('B'), nil},
	}
	for _, test := range tests {
		prepareRepairFromSourceTest(t, test.file, func(ppFilePath string) {
			errors := RepairCharacterLiteralsFromSource(root, ppFilePath)
			if cl.Value != test.expected {
				t.Errorf("RepairCharacterLiteralsFromSource - expected: %x, got: %x", test.expected, cl.Value)
			}
			if test.err != nil && len(errors) == 0 || test.err == nil && len(errors) != 0 {
				t.Errorf("RepairCharacterLiteralsFromSource - error should match: expected: %v, got: %v", test.err, errors)
			} else if test.err != nil && errors[0].Err.Error() != test.err.Error() {
				t.Errorf("RepairCharacterLiteralsFromSource - error should match: expected: %s, got: %s", test.err.Error(), errors[0].Err.Error())
			}
		})
	}
}
