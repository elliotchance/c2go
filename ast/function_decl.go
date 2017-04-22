package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
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

func (n *FunctionDecl) render(program *program.Program) (string, string) {
	out := bytes.NewBuffer([]byte{})

	program.FunctionName = n.Name

	if program.FunctionName == "__istype" ||
		program.FunctionName == "__isctype" ||
		program.FunctionName == "__wcwidth" ||
		program.FunctionName == "__sputc" ||
		program.FunctionName == "__inline_signbitf" ||
		program.FunctionName == "__inline_signbitd" ||
		program.FunctionName == "__inline_signbitl" {
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
		args = append(args, fmt.Sprintf("%s %s", a.Name, types.ResolveType(program, a.Type)))
	}

	if hasBody {
		returnType := getFunctionReturnType(n.Type)

		if program.FunctionName == "main" {
			printLine(out, "func main() {", program.Indent)
		} else {
			printLine(out, fmt.Sprintf("func %s(%s) %s {",
				program.FunctionName, strings.Join(args, ", "),
				types.ResolveType(program, returnType)), program.Indent)
		}

		for _, c := range n.Children {
			if _, ok := c.(*CompoundStmt); ok {
				src, _ := renderExpression(program, c)
				printLine(out, src, program.Indent+1)
			}
		}

		printLine(out, "}\n", program.Indent)

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

func (n *FunctionDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
