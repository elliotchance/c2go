package ast

import "strings"

type TypedefDecl struct {
	Address      string
	Position     string
	Position2    string
	Name         string
	Type         string
	Type2        string
	IsImplicit   bool
	IsReferenced bool
}

func parseTypedefDecl(line string) TypedefDecl {
	groups := groupsFromRegex(
		`<(?P<position>.+?)>(?P<position2> <invalid sloc>| [\w :\d]+?)?(?P<implicit> implicit)?(?P<referenced> referenced)? (?P<name>\w+) '(?P<type>.*?)'(?P<type2>:'.*?')?`,
		line,
	)

	type2 := groups["type2"]
	if type2 != "" {
		type2 = type2[2:len(type2) - 1]
	}

	return TypedefDecl{
		Address: groups["address"],
		Position: groups["position"],
		Position2: strings.TrimSpace(groups["position2"]),
		Name: groups["name"],
		Type: groups["type"],
		Type2: type2,
		IsImplicit: len(groups["implicit"]) > 0,
		IsReferenced: len(groups["referenced"]) > 0,
	}
}
