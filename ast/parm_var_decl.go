package ast

import (
	"strings"
)

type ParmVarDecl struct {
	Address   string
	Position  string
	Position2 string
	Name      string
	Type      string
	Type2     string
	IsUsed    bool
	Children  []Node
}

func parseParmVarDecl(line string) *ParmVarDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<position2> [^ ]+:[\d:]+)?
		(?P<used> used)?
		(?P<name> \w+)?
		 '(?P<type>.*?)'
		(?P<type2>:'.*?')?`,
		line,
	)

	type2 := groups["type2"]
	if type2 != "" {
		type2 = type2[2 : len(type2)-1]
	}

	if strings.Index(groups["position"], "<invalid sloc>") > -1 {
		groups["position"] = "<invalid sloc>"
		groups["position2"] = "<invalid sloc>"
	}

	return &ParmVarDecl{
		Address:   groups["address"],
		Position:  groups["position"],
		Position2: strings.TrimSpace(groups["position2"]),
		Name:      strings.TrimSpace(groups["name"]),
		Type:      groups["type"],
		Type2:     type2,
		IsUsed:    len(groups["used"]) > 0,
		Children:  []Node{},
	}
}

func (n *ParmVarDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
