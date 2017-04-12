package ast

import "strings"

type VarDecl struct {
	Address   string
	Position  string
	Position2 string
	Name      string
	Type      string
	Type2     string
	IsExtern  bool
	IsUsed    bool
	IsCInit   bool
	Children  []interface{}
}

func parseVarDecl(line string) *VarDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>(?P<position2> .+:\d+)?(?P<used> used)?(?P<name> \w+)? '(?P<type>.+?)'(?P<type2>:'.*?')?(?P<extern> extern)?(?P<cinit> cinit)?`,
		line,
	)

	type2 := groups["type2"]
	if type2 != "" {
		type2 = type2[2:len(type2) - 1]
	}

	return &VarDecl{
		Address: groups["address"],
		Position: groups["position"],
		Position2: strings.TrimSpace(groups["position2"]),
		Name: strings.TrimSpace(groups["name"]),
		Type: groups["type"],
		Type2: type2,
		IsExtern: len(groups["extern"]) > 0,
		IsUsed: len(groups["used"]) > 0,
		IsCInit: len(groups["cinit"]) > 0,
		Children: []interface{}{},
	}
}
