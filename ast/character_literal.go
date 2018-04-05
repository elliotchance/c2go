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

		// Use the node position to retrieve the original line from the
		// preprocessed source.
		pos := node.Position()
		line, err :=
			cc.GetLineFromPreprocessedFile(preprocessedFile, pos.File, pos.Line)

		// If there was a problem reading the line we should raise a warning and
		// use the value we have. Hopefully that will be an accurate enough
		// representation.
		if err != nil {
			errs = append(errs, CharacterLiteralError{
				Node: cNode,
				Err:  err,
			})
		}

		// Extract the exact value from the line.
		if pos.Column-1 >= len(line) {
			errs = append(errs, CharacterLiteralError{
				Node: cNode,
				Err:  errors.New("cannot get exact value"),
			})
		} else {
			var i int
			literal := line[pos.Column-1:]
			if _, err := fmt.Sscan(literal, &i); err == nil {
				cNode.Value = i
			} else {
				errs = append(errs, CharacterLiteralError{
					Node: cNode,
					Err:  fmt.Errorf("cannot parse character literal: %v from %s", err, literal),
				})
			}
			fmt.Sscan(line[pos.Column-1:], &cNode.Value)
		}
	}

	return errs
}
