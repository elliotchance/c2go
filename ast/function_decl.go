package ast

import (
	"bytes"
	"fmt"
	"strings"
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
	Children   []Node
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
		Children:   []Node{},
	}
}

func (n *FunctionDecl) render(ast *Ast) (string, string) {
	out := bytes.NewBuffer([]byte{})

	ast.functionName = n.Name

	if ast.functionName == "__istype" ||
		ast.functionName == "__isctype" ||
		ast.functionName == "__wcwidth" ||
		ast.functionName == "__sputc" ||
		ast.functionName == "__inline_signbitf" ||
		ast.functionName == "__inline_signbitd" ||
		ast.functionName == "__inline_signbitl" {
		return "", ""
	}

	hasBody := false
	if len(n.Children) > 0 {
		for _, c := range n.Children {
			if _, ok := c.(*CompoundStmt); ok {
				hasBody = true
			}
		}
	}

	args := []string{}
	for _, a := range getFunctionParams(n) {
		args = append(args, fmt.Sprintf("%s %s", a.Name, resolveType(ast, a.Type)))
	}

	if hasBody {
		returnType := getFunctionReturnType(n.Type)

		if ast.functionName == "main" {
			printLine(out, "func main() {", ast.indent)
		} else {
			printLine(out, fmt.Sprintf("func %s(%s) %s {",
				ast.functionName, strings.Join(args, ", "),
				resolveType(ast, returnType)), ast.indent)
		}

		for _, c := range n.Children {
			if _, ok := c.(*CompoundStmt); ok {
				src, _ := renderExpression(ast, c)
				printLine(out, src, ast.indent+1)
			}
		}

		printLine(out, "}\n", ast.indent)

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

	return out.String(), ""
}
