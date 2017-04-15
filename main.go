package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
)

var (
	printAst = flag.Bool("print-ast", false, "Print AST before translated Go code.")
)

func readAST(data []byte) []string {
	uncolored := regexp.MustCompile(`\x1b\[[\d;]+m`).ReplaceAll(data, []byte{})
	return strings.Split(string(uncolored), "\n")
}

func convertLinesToNodes(lines []string) []interface{} {
	nodes := []interface{}{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// It is tempting to discard null AST nodes, but these may
		// have semantic importance: for example, they represent omitted
		// for-loop conditions, as in for(;;).
		line = strings.Replace(line, "<<<NULL>>>", "NullStmt", 1)

		indentAndType := regexp.MustCompile("^([|\\- `]*)(\\w+)").FindStringSubmatch(line)
		if len(indentAndType) == 0 {
			panic(fmt.Sprintf("Cannot understand line '%s'", line))
		}

		offset := len(indentAndType[1])
		node := Parse(line[offset:])

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

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Start(args []string) string {
	if os.Getenv("GOPATH") == "" {
		panic("The $GOPATH must be set.")
	}

	// 1. Compile it first (checking for errors)
	cFilePath := args[0]

	// 2. Preprocess
	pp, err := exec.Command("clang", "-E", cFilePath).Output()
	Check(err)

	pp_file_path := "/tmp/pp.c"
	err = ioutil.WriteFile(pp_file_path, pp, 0644)
	Check(err)

	// 3. Generate JSON from AST
	ast_pp, err := exec.Command("clang", "-Xclang", "-ast-dump", "-fsyntax-only", pp_file_path).Output()
	Check(err)

	lines := readAST(ast_pp)
	if *printAst {
		for _, l := range lines {
			fmt.Println(l)
		}
		fmt.Println()
	}
	nodes := convertLinesToNodes(lines)
	tree := buildTree(nodes, 0)

	// TODO: allow the user to print the JSON tree:
	//jsonTree := ToJSON(tree)
	//_, err := json.MarshalIndent(jsonTree, " ", "  ")
	//Check(err)

	// 3. Parse C and output Go
	//parts := strings.Split(cFilePath, "/")
	//go_file_path := fmt.Sprintf("%s.go", parts[len(parts) - 1][:len(parts) - 2])
	go_out := bytes.NewBuffer([]byte{})

	Render(go_out, tree[0], "", 0, "")

	// Put together the whole file
	all := "package main\n\nimport (\n"

	for _, importName := range Imports {
		all += fmt.Sprintf("\t\"%s\"\n", importName)
	}

	all += ")\n\n" + go_out.String()

	return all
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: c2go <file.c>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
	}
	fmt.Print(Start(flag.Args()))
}
