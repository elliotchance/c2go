package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
)

type TypedefDecl struct {
	Address      string
	Position     string
	Position2    string
	Name         string
	Type         string
	Type2        string
	IsImplicit   bool
	IsReferenced bool
	Children     []Node
}

func parseTypedefDecl(line string) *TypedefDecl {
	groups := groupsFromRegex(
		`<(?P<position><invalid sloc>|.*?)>
		(?P<position2> <invalid sloc>| col:\d+| line:\d+:\d+)?
		(?P<implicit> implicit)?
		(?P<referenced> referenced)?
		(?P<name> \w+)?
		(?P<type> '.*?')?
		(?P<type2>:'.*?')?`,
		line,
	)

	type2 := groups["type2"]
	if type2 != "" {
		type2 = type2[2 : len(type2)-1]
	}

	return &TypedefDecl{
		Address:      groups["address"],
		Position:     groups["position"],
		Position2:    strings.TrimSpace(groups["position2"]),
		Name:         strings.TrimSpace(groups["name"]),
		Type:         removeQuotes(groups["type"]),
		Type2:        type2,
		IsImplicit:   len(groups["implicit"]) > 0,
		IsReferenced: len(groups["referenced"]) > 0,
		Children:     []Node{},
	}
}

func (n *TypedefDecl) render(program *program.Program) (string, string) {
	out := bytes.NewBuffer([]byte{})
	name := strings.TrimSpace(n.Name)

	if types.TypeIsAlreadyDefined(name) {
		return "", ""
	}

	types.TypeIsNowDefined(name)

	// FIXME: All of the logic here is just to avoid errors, it
	// needs to be fixed up.
	// if ("struct" in node["type"] or "union" in node["type"]) and :
	//     return
	n.Type = strings.Replace(n.Type, "unsigned", "", -1)

	resolvedType := types.ResolveType(program, n.Type)

	if name == "__mbstate_t" {
		program.AddImport("github.com/elliotchance/c2go/darwin")
		resolvedType = "darwin.C__mbstate_t"
	}

	if name == "__darwin_ct_rune_t" {
		program.AddImport("github.com/elliotchance/c2go/darwin")
		resolvedType = "darwin.C__darwin_ct_rune_t"
	}

	if name == "__builtin_va_list" || name == "__qaddr_t" || name == "definition" || name ==
		"_IO_lock_t" || name == "va_list" || name == "fpos_t" || name == "__NSConstantString" || name ==
		"__darwin_va_list" || name == "__fsid_t" || name == "_G_fpos_t" || name == "_G_fpos64_t" {
		return "", ""
	}

	printLine(out, fmt.Sprintf("type %s %s\n", name, resolvedType), program.Indent)

	return out.String(), ""
}

func (n *TypedefDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
