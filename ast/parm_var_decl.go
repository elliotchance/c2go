package ast

import (
	"strings"
)

type ParmVarDecl struct {
	Addr      Address
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
		Addr:      ParseAddress(groups["address"]),
		Position:  groups["position"],
		Position2: strings.TrimSpace(groups["position2"]),
		Name:      strings.TrimSpace(groups["name"]),
		Type:      groups["type"],
		Type2:     type2,
		IsUsed:    len(groups["used"]) > 0,
		Children:  []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ParmVarDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ParmVarDecl) Address() Address {
	return n.Addr
}
