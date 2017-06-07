package ast

// WeakAttr for the WeakAttr node
type WeakAttr struct {
    Address  string
    Position string
    Children []Node
}

func parseWeakAttr(line string) *WeakAttr {
    groups := groupsFromRegex(
        `<(?P<position>.*)>`,
        line,
    )

    return &WeakAttr{
        Address:  groups["address"],
        Position: groups["position"],
        Children: []Node{},
    }
}

// AddChild method to implements Node interface
func (a *WeakAttr) AddChild(node Node) {
    a.Children = append(a.Children, node)
}
