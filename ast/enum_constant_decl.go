package ast

import (
	"fmt"

	"github.com/elliotchance/c2go/program"
)

type EnumConstantDecl struct {
	Address    string
	Position   string
	Position2  string
	Referenced bool
	Name       string
	Type       string
	Children   []Node
}

func parseEnumConstantDecl(line string) *EnumConstantDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		( (?P<position2>[^ ]+))?
		( (?P<referenced>referenced))?
		 (?P<name>.+)
		 '(?P<type>.+?)'`,
		line,
	)

	return &EnumConstantDecl{
		Address:    groups["address"],
		Position:   groups["position"],
		Position2:  groups["position2"],
		Referenced: len(groups["referenced"]) > 0,
		Name:       groups["name"],
		Type:       groups["type"],
		Children:   []Node{},
	}
}

func (n *EnumConstantDecl) render(program *program.Program) (string, string) {
	value := "iota"

	// Special cases for linux ctype.h
	switch n.Name {
	case "_ISupper":
		value = "((1 << (0)) << 8)"
	case "_ISlower":
		value = "((1 << (1)) << 8)"
	case "_ISalpha":
		value = "((1 << (2)) << 8)"
	case "_ISdigit":
		value = "((1 << (3)) << 8)"
	case "_ISxdigit":
		value = "((1 << (4)) << 8)"
	case "_ISspace":
		value = "((1 << (5)) << 8)"
	case "_ISprint":
		value = "((1 << (6)) << 8)"
	case "_ISgraph":
		value = "((1 << (7)) << 8)"
	case "_ISblank":
		value = "((1 << (8)) >> 8)"
	case "_IScntrl":
		value = "((1 << (9)) >> 8)"
	case "_ISpunct":
		value = "((1 << (10)) >> 8)"
	case "_ISalnum":
		value = "((1 << (11)) >> 8)"
	default:
		if len(n.Children) > 0 {
			value, _ = n.Children[0].render(program)
		}
	}

	line := fmt.Sprintf("%s = %s", n.Name, value)

	return line, n.Type
}

func (n *EnumConstantDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
