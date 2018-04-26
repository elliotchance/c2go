package ast

import (
	"errors"
	"fmt"
	"github.com/elliotchance/c2go/cc"
	"github.com/elliotchance/c2go/util"
	"reflect"
)

// CharacterLiteral is type of character literal
type CharacterLiteral struct {
	Addr       Address
	Pos        Position
	Type       string
	Value      int
	ChildNodes []Node
}

func parseCharacterLiteral(line string) *CharacterLiteral {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<value>\\d+)",
		line,
	)

	return &CharacterLiteral{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Value:      util.Atoi(groups["value"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CharacterLiteral) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CharacterLiteral) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CharacterLiteral) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CharacterLiteral) Position() Position {
	return n.Pos
}

// CharacterLiteralError represents one instance of an error where the exact
// character value of a CharacterLiteral could not be determined from the
// original source. See RepairCharacterLiteralsFromSource for a full explanation.
type CharacterLiteralError struct {
	Node *CharacterLiteral
	Err  error
}

// RepairCharacterLiteralsFromSource finds the exact values of character literals
// by reading their values directly from the preprocessed source.
//
// This is to solve issue #663, sometime clang serializes hex encoded character literals
// as a weird number, e.g. '\xa0' => 4294967200
//
// The only solution is to read the original character literal from the source
// code. We can do this by using the positional information on the node.
//
// If the character literal cannot be resolved for any reason the original value
// will remain. This function will return all errors encountered.
func RepairCharacterLiteralsFromSource(rootNode Node, preprocessedFile string) []CharacterLiteralError {
	errs := []CharacterLiteralError{}
	characterLiteralNodes :=
		GetAllNodesOfType(rootNode, reflect.TypeOf((*CharacterLiteral)(nil)))

	for _, node := range characterLiteralNodes {
		cNode := node.(*CharacterLiteral)
		pos := node.Position()
		var (
			err      error
			lastLine = pos.LineEnd
			i        int
		)
		if lastLine == 0 {
			lastLine = pos.Line
		}
		for line := pos.Line; line <= lastLine; line++ {
			i, err = parseCharacterLiteralFromPosition(preprocessedFile, pos, line)
			if err == nil {
				cNode.Value = i
				break
			}
		}
		if err != nil {
			errs = append(errs, CharacterLiteralError{
				Node: cNode,
				Err:  err,
			})
		}
	}

	return errs
}

func parseCharacterLiteralFromPosition(preprocessedFile string, pos Position, lineNbr int) (ret int, err error) {
	// Use the node position to retrieve the original line from the
	// preprocessed source.
	line, err :=
		cc.GetLineFromPreprocessedFile(preprocessedFile, pos.File, lineNbr)

	// If there was a problem reading the line we should raise a warning and
	// use the value we have. Hopefully that will be an accurate enough
	// representation.
	if err != nil {
		return 0, err
	}

	// Extract the exact value from the line.
	if pos.Column-1 >= len(line) {
		return 0, errors.New("cannot get exact value")
	}
	literal := line[pos.Column-1:]
	if ret, err = parseCharacterLiteralFromSource(literal); err == nil {
		return ret, nil
	}
	return 0, fmt.Errorf("cannot parse character literal: %v from %s", err, literal)
}

func parseCharacterLiteralFromSource(literal string) (ret int, err error) {
	runes := []rune(literal)
	if len(runes) < 1 {
		return 0, fmt.Errorf("character literal to short")
	}
	// Consume leading '
	switch runes[0] {
	case '\'':
		runes = runes[1:]
	case 'u', 'U', 'L':
		if len(runes) < 2 {
			return 0, fmt.Errorf("character literal to short")
		}
		if runes[1] == '\'' {
			runes = runes[2:]
		} else if runes[1] == '8' {
			if len(runes) < 3 {
				return 0, fmt.Errorf("character literal to short")
			} else if runes[2] != '\'' {
				return 0, fmt.Errorf("illegal character '%s' at index 2", string(runes[2]))
			}
			runes = runes[3:]
		} else {
			return 0, fmt.Errorf("illegal character '%s' at index 1", string(runes[1]))
		}
	default:
		return 0, fmt.Errorf("illegal character '%s' at index 0", string(runes[0]))
	}

	// we need place for at least 1 character and '
	if len(runes) < 1 {
		return 0, fmt.Errorf("unexpected end of character literal")
	}
	// decode character literal
	var r rune
	var i int
	switch runes[0] {
	case '\'':
		return 0, fmt.Errorf("empty character literal")
	case '\\':
		if len(runes) < 2 {
			return 0, fmt.Errorf("unexpected end of character literal")
		}
		r, i, err = decodeEscapeSequence(runes)
	default:
		r = runes[0]
		i = 1
	}
	if err != nil {
		return 0, err
	}
	if len(runes) <= i {
		return 0, fmt.Errorf("unexpected end of character literal")
	}
	if runes[i] != '\'' {
		return 0, fmt.Errorf("does not support multi-character literals")
	}
	return int(r), nil
}

// escape-sequence		{simple-sequence}|{octal-escape-sequence}|{hexadecimal-escape-sequence}|{universal-character-name}
// simple-sequence		\\['\x22?\\abfnrtv]
// octal-escape-sequence	\\{octal-digit}{octal-digit}?{octal-digit}?
// hexadecimal-escape-sequence	\\x{hexadecimal-digit}+
func decodeEscapeSequence(runes []rune) (rune, int, error) {
	if runes[0] != '\\' {
		panic("internal error")
	}

	r := runes[1]
	switch r {
	case '\'', '"', '?', '\\':
		return r, 2, nil
	case 'a':
		return 7, 2, nil
	case 'b':
		return 8, 2, nil
	case 'f':
		return 12, 2, nil
	case 'n':
		return 10, 2, nil
	case 'r':
		return 13, 2, nil
	case 't':
		return 9, 2, nil
	case 'v':
		return 11, 2, nil
	case 'x':
		v, n := 0, 2
	loop2:
		for _, r := range runes[2:] {
			switch {
			case r >= '0' && r <= '9', r >= 'a' && r <= 'f', r >= 'A' && r <= 'F':
				v = v<<4 | decodeHex(r)
				n++
				if n >= 4 {
					break loop2
				}
			default:
				break loop2
			}
		}
		return rune(v & 0xff), n, nil
	case 'u', 'U':
		v, n := decodeUCN(runes)
		return v, n, nil
	}

	if r < '0' || r > '7' {
		return 0, 0, fmt.Errorf("illegal character '%s'", string(r))
	}

	v, n := 0, 1
loop:
	for _, r := range runes[1:] {
		switch {
		case r >= '0' && r <= '7':
			v = v<<3 | (int(r) - '0')
			n++
			if n >= 4 {
				break loop
			}
		default:
			break loop
		}
	}
	return rune(v), n, nil
}

func decodeHex(r rune) int {
	switch {
	case r >= '0' && r <= '9':
		return int(r) - '0'
	case r >= 'a' && r <= 'f', r >= 'A' && r <= 'F':
		x := int(r) &^ 0x20
		return x - 'A' + 10
	default:
		return -1
	}
}

// universal-character-name	\\u{hex-quad}|\\U{hex-quad}{hex-quad}
func decodeUCN(runes []rune) (rune, int) {
	if runes[0] != '\\' {
		panic("internal error")
	}

	runes = runes[1:]
	switch runes[0] {
	case 'u':
		hq, n := decodeHexQuad(runes[1:])
		return rune(hq), n + 2
	case 'U':
		hq, n := decodeHexQuad(runes[1:])
		if n == 4 {
			hq2, n2 := decodeHexQuad(runes[5:])
			hq = hq << (4 * uint(n2))
			hq = hq | hq2
			n = n + n2
		}
		return rune(hq), n + 2
	default:
		panic("internal error")
	}
}

// hex-quad	{hexadecimal-digit}{hexadecimal-digit}{hexadecimal-digit}{hexadecimal-digit}
func decodeHexQuad(runes []rune) (int, int) {
	v, n := 0, 0
	for _, r := range runes[:4] {
		h := decodeHex(r)
		if h < 0 {
			break
		}
		v = v<<4 | h
		n++
	}
	return v, n
}
