package main

import (
	"strings"
	"fmt"
	"bytes"
)

type FunctionDecl struct {
	Address    string
	Position   string
	Prev       string
	Position2  string
	Name       string
	Type       string
	IsExtern   bool
	IsImplicit bool
	IsUsed     bool
	Children   []interface{}
}

func parseFunctionDecl(line string) *FunctionDecl {
	groups := groupsFromRegex(
		`(?P<prev>prev [0-9a-fx]+ )?
		<(?P<position1>.*?)>
		(?P<position2> <scratch space>[^ ]+| [^ ]+)?
		(?P<implicit> implicit)?
		(?P<used> used)?
		 (?P<name>[_\w]+)
		 '(?P<type>.*)
		'(?P<extern> extern)?`,
		line,
	)

	prev := groups["prev"]
	if prev != "" {
		prev = prev[5 : len(prev)-1]
	}

	return &FunctionDecl{
		Address:    groups["address"],
		Position:   groups["position1"],
		Prev:       prev,
		Position2:  strings.TrimSpace(groups["position2"]),
		Name:       groups["name"],
		Type:       groups["type"],
		IsExtern:   len(groups["extern"]) > 0,
		IsImplicit: len(groups["implicit"]) > 0,
		IsUsed:     len(groups["used"]) > 0,
		Children:   []interface{}{},
	}
}

func (n *FunctionDecl) RenderLine(out *bytes.Buffer, functionName string, indent int, returnType string) {
	functionName = strings.TrimSpace(n.Name)

	if functionName == "__istype" || functionName == "__isctype" ||
		functionName == "__wcwidth" || functionName == "__sputc" ||
		functionName == "__inline_signbitf" ||
		functionName == "__inline_signbitd" ||
		functionName == "__inline_signbitl" {
		return
	}

	has_body := false
	if len(n.Children) > 0 {
		for _, c := range n.Children {
			if _, ok := c.(*CompoundStmt); ok {
				has_body = true
			}
		}
	}

	args := []string{}
	for _, a := range getFunctionParams(n) {
		args = append(args, fmt.Sprintf("%s %s", a.Name, resolveType(a.Type)))
	}

	if has_body {
		returnType := getFunctionReturnType(n.Type)

		if functionName == "main" {
			printLine(out, "func main() {", indent)
		} else {
			printLine(out, fmt.Sprintf("func %s(%s) %s {",
				functionName, strings.Join(args, ", "),
				resolveType(returnType)), indent)
		}

		for _, c := range n.Children {
			if _, ok := c.(*CompoundStmt); ok {
				Render(out, c, functionName,
					indent+1, n.Type)
			}
		}

		printLine(out, "}\n", indent)

		params := []string{}
		for _, v := range getFunctionParams(n) {
			params = append(params, v.Type)
		}

		addFunctionDefinition(FunctionDefinition{
			Name:          n.Name,
			ReturnType:    getFunctionReturnType(n.Type),
			ArgumentTypes: params,
		})
	}
}
