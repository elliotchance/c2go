package ast

import "github.com/elliotchance/c2go/program"

type WarnUnusedResultAttr struct {
	Address  string
	Position string
	Children []interface{}
}

func parseWarnUnusedResultAttr(line string) *WarnUnusedResultAttr {
	groups := groupsFromRegex(`<(?P<position>.*)> warn_unused_result`, line)

	return &WarnUnusedResultAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *WarnUnusedResultAttr) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *WarnUnusedResultAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
