package ast

import (
	"strings"
)

type VarDecl struct {
	Addr         Address
	Position     string
	Position2    string
	Name         string
	Type         string
	Type2        string
	IsExtern     bool
	IsUsed       bool
	IsCInit      bool
	IsReferenced bool
	Children     []Node
}

func parseVarDecl(line string) *VarDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>(?P<position2> .+:\d+)?
		(?P<used> used)?
		(?P<referenced> referenced)?
		(?P<name> \w+)?
		 '(?P<type>.+?)'
		(?P<type2>:'.*?')?
		(?P<extern> extern)?
		(?P<cinit> cinit)?`,
		line,
	)

	type2 := groups["type2"]
	if type2 != "" {
		type2 = type2[2 : len(type2)-1]
	}

	return &VarDecl{
		Addr:         ParseAddress(groups["address"]),
		Position:     groups["position"],
		Position2:    strings.TrimSpace(groups["position2"]),
		Name:         strings.TrimSpace(groups["name"]),
		Type:         groups["type"],
		Type2:        type2,
		IsExtern:     len(groups["extern"]) > 0,
		IsUsed:       len(groups["used"]) > 0,
		IsCInit:      len(groups["cinit"]) > 0,
		IsReferenced: len(groups["referenced"]) > 0,
		Children:     []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *VarDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *VarDecl) Address() Address {
	return n.Addr
}
