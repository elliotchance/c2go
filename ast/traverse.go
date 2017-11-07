package ast

import (
	"reflect"
)

// GetAllNodesOfType returns all of the nodes of the tree that match the type
// provided. The type should be a pointer to an object in the ast package.
//
// The nodes returned may reference each other and there is no guaranteed order
// in which the nodes are returned.
func GetAllNodesOfType(root Node, t reflect.Type) []Node {
	nodes := []Node{}

	if root == nil {
		return []Node{}
	}

	if reflect.TypeOf(root) == t {
		nodes = append(nodes, root)
	}

	for _, c := range root.Children() {
		nodes = append(nodes, GetAllNodesOfType(c, t)...)
	}

	return nodes
}
