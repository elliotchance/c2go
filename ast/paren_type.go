package ast

import (
	"github.com/elliotchance/c2go/program"
)

type ParenType struct {
	Address  string
	Type     string
	Sugar    bool
	Children []interface{}
}

func parseParenType(line string) *ParenType {
	groups := groupsFromRegex(`'(?P<type>.*?)' sugar`, line)

	return &ParenType{
		Address:  groups["address"],
		Type:     groups["type"],
		Sugar:    true,
		Children: []interface{}{},
	}
}

func (n *ParenType) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *ParenType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
