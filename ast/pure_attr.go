package ast

import "github.com/elliotchance/c2go/program"

type PureAttr struct {
	Address   string
	Position  string
	Implicit  bool
	Inherited bool
	Children  []interface{}
}

func parsePureAttr(line string) *PureAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<inherited> Inherited)?
		(?P<implicit> Implicit)?`,
		line,
	)

	return &PureAttr{
		Address:   groups["address"],
		Position:  groups["position"],
		Implicit:  len(groups["implicit"]) > 0,
		Inherited: len(groups["inherited"]) > 0,
		Children:  []interface{}{},
	}
}

func (n *PureAttr) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *PureAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
