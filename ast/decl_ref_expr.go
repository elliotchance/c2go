package ast

type DeclRefExpr struct {
	Address  string
	Position string
	Type     string
	Lvalue   bool
	For      string
	Address2 string
	Name     string
	Type2    string
	Children []Node
}

func parseDeclRefExpr(line string) *DeclRefExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'
		.*?
		(?P<lvalue> lvalue)?
		 (?P<for>\w+)
		 (?P<address2>[0-9a-fx]+)
		 '(?P<name>.*?)'
		 '(?P<type2>.*?)'`,
		line,
	)

	return &DeclRefExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Lvalue:   len(groups["lvalue"]) > 0,
		For:      groups["for"],
		Address2: groups["address2"],
		Name:     groups["name"],
		Type2:    groups["type2"],
		Children: []Node{},
	}
}

func (n *DeclRefExpr) render(ast *Ast) (string, string) {
	name := n.Name

	if name == "argc" {
		name = "len(os.Args)"
		ast.addImport("os")
	} else if name == "argv" {
		name = "os.Args"
		ast.addImport("os")
	}

	return name, n.Type
}

func (n *DeclRefExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
