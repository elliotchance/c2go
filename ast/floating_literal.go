package ast

import (
	"fmt"

	"github.com/elliotchance/c2go/program"
)

type FloatingLiteral struct {
	Address  string
	Position string
	Type     string
	Value    float64
	Children []Node
}

func parseFloatingLiteral(line string) *FloatingLiteral {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<value>.+)",
		line,
	)

	return &FloatingLiteral{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Value:    atof(groups["value"]),
		Children: []Node{},
	}
}

func (n *FloatingLiteral) render(program *program.Program) (string, string) {
	return fmt.Sprintf("%f", n.Value), "double"
}

func (n *FloatingLiteral) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
