package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func removeQuotes(s string) string {
	s = strings.TrimSpace(s)

	if s == `""` {
		return ""
	}
	if s == `''` {
		return ""
	}

	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-2]
	}
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1]
	}

	return s
}

func atof(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}

	return f
}

// Atos - ASTree to string
// Typically using for debug
func Atos(node Node) string {
	j, err := json.Marshal(node)
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	err = json.Indent(&out, j, "", "\t")
	if err != nil {
		panic(err)
	}
	var str string
	str += fmt.Sprint("==== START OF AST tree ====\n")
	str += out.String()
	str += TypesTree(node)
	str += fmt.Sprint("==== END OF AST tree ====\n")
	return str
}

// TypesTree - return tree of types for AST node
func TypesTree(node Node) (str string) {
	str += fmt.Sprintf("\nTypes tree:\n")
	str += typesTree(node, 0)
	return str
}

func typesTree(node Node, depth int) (str string) {
	if node == (Node)(nil) {
		return ""
	}
	for i := 0; i < depth; i++ {
		str += "\t"
	}
	str += fmt.Sprintf("%T\n", node)
	depth++
	if len(node.Children()) > 0 {
		for _, n := range node.Children() {
			str += typesTree(n, depth)
		}
	}
	return str
}
