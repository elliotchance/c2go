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

func (n *FunctionDecl) render(p *program.Program) (string, string) {
	out := bytes.NewBuffer([]byte{})
	p.FunctionName = n.Name

	// Always register the new function. Only from this point onwards will
	// we be allowed to refer to the function.
	if program.GetFunctionDefinition(p.FunctionName) == nil {
		program.AddFunctionDefinition(program.FunctionDefinition{
			Name:       n.Name,
			ReturnType: getFunctionReturnType(n.Type),
			// FIXME
			ArgumentTypes: []string{},
			Substitution:  "",
		})
	}

	// If the function has a direct substitute in Go we do not want to
	// output the C definition of it.
	if f := program.GetFunctionDefinition(p.FunctionName); f != nil &&
		f.Substitution != "" {
		return "", ""
	}

	if p.FunctionName == "__istype" ||
		p.FunctionName == "__isctype" ||
		p.FunctionName == "__wcwidth" ||
		p.FunctionName == "__sputc" ||
		p.FunctionName == "__inline_signbitf" ||
		p.FunctionName == "__inline_signbitd" ||
		p.FunctionName == "__inline_signbitl" {
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
		args = append(args, fmt.Sprintf("%s %s", a.Name, types.ResolveType(p, a.Type)))
	}

	if hasBody {
		returnType := getFunctionReturnType(n.Type)

		if p.FunctionName == "main" {
			printLine(out, "func main() {", p.Indent)
		} else {
			printLine(out, fmt.Sprintf("func %s(%s) %s {",
				p.FunctionName, strings.Join(args, ", "),
				types.ResolveType(p, returnType)), p.Indent)
		}

		for _, c := range n.Children {
			if _, ok := c.(*CompoundStmt); ok {
				src, _ := renderExpression(p, c)
				printLine(out, src, p.Indent+1)
			}
		}

		printLine(out, "}\n", p.Indent)

		params := []string{}
		for _, v := range getFunctionParams(n) {
			params = append(params, v.Type)
		}

		program.AddFunctionDefinition(program.FunctionDefinition{
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
