package ast

type UnaryOperator struct {
	Addr       Address
	Position   string
	Type       string
	IsLvalue   bool
	IsPrefix   bool
	Operator   string
	ChildNodes []Node
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
		Addr:       ParseAddress(groups["address"]),
		Position:   groups["position"],
		Type:       groups["type"],
		IsLvalue:   len(groups["lvalue"]) > 0,
		IsPrefix:   len(groups["prefix"]) > 0,
		Operator:   groups["operator"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *UnaryOperator) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *UnaryOperator) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *UnaryOperator) Children() []Node {
	return n.ChildNodes
}
