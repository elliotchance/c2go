package ast

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/elliotchance/c2go/cc"
)

// FloatingLiteral is type of float literal
type FloatingLiteral struct {
	Addr       Address
	Pos        Position
	Type       string
	Value      float64
	ChildNodes []Node
}

func parseFloatingLiteral(line string) *FloatingLiteral {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<value>.+)",
		line,
	)

	return &FloatingLiteral{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Value:      atof(groups["value"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FloatingLiteral) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *FloatingLiteral) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *FloatingLiteral) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *FloatingLiteral) Position() Position {
	return n.Pos
}

// FloatingLiteralError represents one instance of an error where the exact
// floating point value of a FloatingLiteral could not be determined from the
// original source. See RepairFloatingLiteralsFromSource for a full explanation.
type FloatingLiteralError struct {
	Node *FloatingLiteral
	Err  error
}

// RepairFloatingLiteralsFromSource finds the exact values of floating literals
// by reading their values directly from the preprocessed source.
//
// The clang AST only serializes floating point values in scientific notation
// with 7 significant digits. This is not enough when dealing with precise
// numbers.
//
// The only solution is to read the original floating literal from the source
// code. We can do this by using the positional information on the node.
//
// If the floating literal cannot be resolved for any reason the original value
// will remain. This function will return all errors encountered.
func RepairFloatingLiteralsFromSource(rootNode Node, preprocessedFile string) []FloatingLiteralError {
	errs := []FloatingLiteralError{}
	floatingLiteralNodes :=
		GetAllNodesOfType(rootNode, reflect.TypeOf((*FloatingLiteral)(nil)))

	for _, node := range floatingLiteralNodes {
		fNode := node.(*FloatingLiteral)

		// Use the node position to retrieve the original line from the
		// preprocessed source.
		pos := node.Position()
		line, err :=
			cc.GetLineFromPreprocessedFile(preprocessedFile, pos.File, pos.Line)

		// If there was a problem reading the line we should raise a warning and
		// use the value we have. Hopefully that will be an accurate enough
		// representation.
		if err != nil {
			errs = append(errs, FloatingLiteralError{
				Node: fNode,
				Err:  err,
			})
		}

		// Extract the exact value from the line.
		if pos.Column-1 >= len(line) {
			errs = append(errs, FloatingLiteralError{
				Node: fNode,
				Err:  errors.New("cannot get exact value"),
			})
		} else {
			var f float64
			literal := line[pos.Column-1:]
			if _, err := fmt.Sscan(literal, &f); err == nil {
				fNode.Value = f
			} else {
				errs = append(errs, FloatingLiteralError{
					Node: fNode,
					Err:  fmt.Errorf("cannot parse float: %v from %s", err, literal),
				})
			}
		}
	}

	return errs
}
