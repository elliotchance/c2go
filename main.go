package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/elliotchance/c2go/ast"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func readAST() []string {
	reader := bufio.NewReader(os.Stdin)
	data := []byte{}

	for {
		buf := make([]byte, 16384)
		bytesRead, err := reader.Read(buf)
		if err != nil && err.Error() != "EOF" {
			panic(err)
		}
		if bytesRead == 0 {
			break
		}
		data = append(data, buf[0:bytesRead]...)
	}

	uncolored := regexp.MustCompile(`\x1b\[[\d;]+m`).ReplaceAll(data, []byte{})
	return strings.Split(string(uncolored), "\n")
}

func convertLinesToNodes(lines []string) []interface{} {
	nodes := []interface{}{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// This will need to be handled more gracefully... I'm not even
		// sure what this means?
		if strings.Index(line, "<<<NULL>>>") >= 0 {
			continue
		}

		indentAndType := regexp.MustCompile("^([|\\- `]*)(\\w+)").FindStringSubmatch(line)
		if len(indentAndType) == 0 {
			panic(fmt.Sprintf("Can not understand line '%s'", line))
		}

		//nodeType := indentAndType[2]
		offset := len(indentAndType[1])
		//try:
		node := ast.Parse(line[offset:])
		//except KeyError:
		//    print("There is no regex for '%s'." % node_type)
		//    print("I will print out all the lines so a regex can be created:\n")
		//
		//    for line in lines:
		//        //s = re.search(r'^([|\- `]*)(\w+)', line)
		//        if s is not None and node_type == s.group(2):
		//            print(line[offset:])
		//
		//    sys.exit(1)

		//if result is None:
		//    print("Can not understand line '%s'" % line)
		//    sys.exit(1)

		//node = result.groupdict()

		//node['node'] = node_type

		indentLevel := len(indentAndType[1]) / 2
		nodes = append(nodes, []interface{}{indentLevel, node})
	}

	return nodes
}

// buildTree convert an array of nodes, each prefixed with a depth into a tree.
func buildTree(nodes []interface{}, depth int) []interface{} {
	if len(nodes) == 0 {
		return []interface{}{}
	}

	// Split the list into sections, treat each section as a a tree with its own root.
	sections := [][]interface{}{}
	for _, node := range nodes {
		if node.([]interface{})[0] == depth {
			sections = append(sections, []interface{}{node})
		} else {
			sections[len(sections)-1] = append(sections[len(sections)-1], node)
		}
	}

	results := []interface{}{}
	for _, section := range sections {
		slice := []interface{}{}
		for _, n := range section {
			if n.([]interface{})[0].(int) > depth {
				slice = append(slice, n)
			}
		}

		children := buildTree(slice, depth+1)
		result := section[0].([]interface{})[1]

		if len(children) > 0 {
			c := reflect.ValueOf(result).Elem().FieldByName("Children")
			slice := reflect.AppendSlice(c, reflect.ValueOf(children))
			c.Set(slice)
		}

		results = append(results, result)
	}

	return results
}

func ToJSON(tree []interface{}) []map[string]interface{} {
	r := make([]map[string]interface{}, len(tree))

	for j, n := range tree {
		rn := reflect.ValueOf(n).Elem()
		r[j] = make(map[string]interface{})
		r[j]["node"] = rn.Type().Name()

		for i := 0; i < rn.NumField(); i++ {
			name := strings.ToLower(rn.Type().Field(i).Name)
			value := rn.Field(i).Interface()

			if name == "children" {
				v := value.([]interface{})

				if len(v) == 0 {
					continue
				}

				value = ToJSON(v)
			}

			r[j][name] = value
		}
	}

	return r
}

func main() {
	lines := readAST()
	nodes := convertLinesToNodes(lines)
	tree := buildTree(nodes, 0)
	jsonTree := ToJSON(tree)

	out, err := json.MarshalIndent(jsonTree, " ", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Print(string(out))
}
