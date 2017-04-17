package ast

import (
	"fmt"
	"strings"
)

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
	Children  []Node
}

func parseVarDecl(line string) *VarDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>(?P<position2> .+:\d+)?
		(?P<used> used)?
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
		Address:   groups["address"],
		Position:  groups["position"],
		Position2: strings.TrimSpace(groups["position2"]),
		Name:      strings.TrimSpace(groups["name"]),
		Type:      groups["type"],
		Type2:     type2,
		IsExtern:  len(groups["extern"]) > 0,
		IsUsed:    len(groups["used"]) > 0,
		IsCInit:   len(groups["cinit"]) > 0,
		Children:  []Node{},
	}
}

func (n *VarDecl) render(ast *Ast) (string, string) {
	theType := resolveType(ast, n.Type)
	name := n.Name

	// FIXME: These names dont seem to work when testing more than 1 file
	if name == "_LIB_VERSION" ||
		name == "_IO_2_1_stdin_" ||
		name == "_IO_2_1_stdout_" ||
		name == "_IO_2_1_stderr_" ||
		name == "stdin" ||
		name == "stdout" ||
		name == "stderr" {
		return "", ""
	}

	// Go does not allow the name of a variable to be called "type".
	// For the moment I will rename this to avoid the error.
	if name == "type" {
		name = "type_"
	}

	suffix := ""
	if len(n.Children) > 0 {
		children := n.Children
		src, _ := renderExpression(ast, children[0])
		suffix = fmt.Sprintf(" = %s", src)
	}

	if suffix == " = (0)" {
		suffix = " = nil"
	}

	return fmt.Sprintf("var %s %s%s", name, theType, suffix), "unknown3"
}

func (n *VarDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
