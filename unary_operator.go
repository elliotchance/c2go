package main

import "fmt"

type UnaryOperator struct {
	Address  string
	Position string
	Type     string
	IsLvalue bool
	IsPrefix bool
	Operator string
	Children []interface{}
}

func parseUnaryOperator(line string) *UnaryOperator {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'
		(?P<lvalue> lvalue)?
		(?P<prefix> prefix)?
		(?P<postfix> postfix)?
		 '(?P<operator>.*?)'`,
		line,
	)

	return &UnaryOperator{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		IsLvalue: len(groups["lvalue"]) > 0,
		IsPrefix: len(groups["prefix"]) > 0,
		Operator: groups["operator"],
		Children: []interface{}{},
	}
}

func (n *UnaryOperator) Render() []string {
	operator := n.Operator
	expr := renderExpression(n.Children[0])

	if operator == "!" {
		if expr[1] == "bool" {
			return []string{fmt.Sprintf("!(%s)", expr[0]), expr[1]}
		}

		addImport("github.com/elliotchance/c2go/noarch")

		functionName := fmt.Sprintf("noarch.Not%s", ucfirst(expr[1]))
		return []string{
			fmt.Sprintf("%s(%s)", functionName, expr[0]),
			expr[1],
		}
	}

	if operator == "*" {
		if expr[1] == "const char *" {
			return []string{fmt.Sprintf("%s[0]", expr[0]), "char"}
		}

		return []string{fmt.Sprintf("*%s", expr[0]), "int"}
	}

	if operator == "++" {
		return []string{fmt.Sprintf("%s += 1", expr[0]), expr[1]}
	}

	if operator == "~" {
		operator = "^"
	}

	return []string{fmt.Sprintf("%s%s", operator, expr[0]), expr[1]}
}
