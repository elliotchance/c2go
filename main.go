package main

import (
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
)

var (
	printAst = flag.Bool("print-ast", false, "Print AST before translated Go code.")
)

func readAST(data []byte) []string {
	uncolored := regexp.MustCompile(`\x1b\[[\d;]+m`).ReplaceAll(data, []byte{})
	return strings.Split(string(uncolored), "\n")
}

type treeNode struct {
	indent int
	node   ast.Node
}

func convertLinesToNodes(lines []string) []treeNode {
	nodes := []treeNode{}
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
		node := ast.Parse(line[offset:])

		indentLevel := len(indentAndType[1]) / 2
		nodes = append(nodes, treeNode{indentLevel, node})
	}

	return nodes
}

// buildTree convert an array of nodes, each prefixed with a depth into a tree.
func buildTree(nodes []treeNode, depth int) []ast.Node {
	if len(nodes) == 0 {
		return []ast.Node{}
	}

	// Split the list into sections, treat each section as a a tree with its own root.
	sections := [][]treeNode{}
	for _, node := range nodes {
		if node.indent == depth {
			sections = append(sections, []treeNode{node})
		} else {
			sections[len(sections)-1] = append(sections[len(sections)-1], node)
		}
	}

	results := []ast.Node{}
	for _, section := range sections {
		slice := []treeNode{}
		for _, n := range section {
			if n.indent > depth {
				slice = append(slice, n)
			}
		}

		children := buildTree(slice, depth+1)
		for _, child := range children {
			section[0].node.AddChild(child)
		}
		results = append(results, section[0].node)
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

	_, err := os.Stat(cFilePath)
	Check(err)

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

	// Render(go_out, tree[0], "", 0, "")
	p := program.NewProgram()
	goOut := ast.Render(p, tree[0].(ast.Node))

	// Format the code
	goOutFmt, err := format.Source([]byte(goOut))
	if err != nil {
		panic(err.Error() + "\n\n" + goOut)
	}

	// Put together the whole file
	all := "package main\n\nimport (\n"

	for _, importName := range p.Imports() {
		all += fmt.Sprintf("\t\"%s\"\n", importName)
	}

	all += ")\n\n" + string(goOutFmt)

	return all
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file.c>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Print(Start(flag.Args()))
}
