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
	name := n.Name

	if program.TypeIsAlreadyDefined(name) {
		return "", ""
	}

	program.TypeIsNowDefined(name)

	resolvedType := types.ResolveType(program, n.Type)
	// There is a case where the name of the type is also the definition,
	// like:
	//
	//     type _RuneEntry _RuneEntry
	//
	// This of course is impossible and will cause the Go not to compile.
	// It itself is caused by lack of understanding (at this time) about
	// certain scenarios that types are defined as. The above example comes
	// from:
	//
	//     typedef struct {
	//        // ... some fields
	//     } _RuneEntry;
	//
	// Until which time that we actually need this to work I am going to
	// suppress these.
	if name == resolvedType {
		return "", ""
	}

	if name == "__mbstate_t" {
		resolvedType = program.ImportType("github.com/elliotchance/c2go/darwin.C__mbstate_t")
	}

	if name == "__darwin_ct_rune_t" {
		resolvedType = program.ImportType("github.com/elliotchance/c2go/darwin.Darwin_ct_rune_t")
	}

	// A bunch of random stuff to ignore... I really should deal with these.
	if name == "__builtin_va_list" ||
		name == "__qaddr_t" ||
		name == "definition" ||
		name == "_IO_lock_t" ||
		name == "va_list" ||
		name == "fpos_t" ||
		name == "__NSConstantString" ||
		name == "__darwin_va_list" ||
		name == "__fsid_t" ||
		name == "_G_fpos_t" ||
		name == "_G_fpos64_t" || name == "__locale_t" || name == "locale_t" {
		return "", ""
	}

	printLine(out, fmt.Sprintf("type %s %s\n", name, resolvedType), program.Indent)

	return out.String(), ""
}

func (n *TypedefDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
