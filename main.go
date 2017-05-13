package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/transpiler"
)

const Version = "0.10.5"

var (
	printAst = flag.Bool("print-ast", false, "Print AST before translated Go code.")
	version  = flag.Bool("version", false, "Print the version and exit.")
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

func Check(prefix string, e error) {
	if e != nil {
		panic(prefix + e.Error())
	}
}

func Start(args []string) string {
	if os.Getenv("GOPATH") == "" {
		panic("The $GOPATH must be set.")
	}

	// 1. Compile it first (checking for errors)
	cFilePath := args[0]

	_, err := os.Stat(cFilePath)
	Check("", err)

	// 2. Preprocess
	pp, err := exec.Command("clang", "-E", cFilePath).Output()
	Check("preprocess failed: ", err)

	pp_file_path := path.Join(os.TempDir(), "pp.c")
	err = ioutil.WriteFile(pp_file_path, pp, 0644)
	Check("writing to /tmp/pp.c failed: ", err)

	// 3. Generate JSON from AST
	ast_pp, err := exec.Command("clang", "-Xclang", "-ast-dump", "-fsyntax-only", pp_file_path).Output()
	if err != nil {
		// If clang fails it still prints out the AST, so we have to run it
		// again to get the real error.
		errBody, _ := exec.Command("clang", pp_file_path).CombinedOutput()

		panic("clang failed: " + err.Error() + ":\n\n" + string(errBody))
	}

	lines := readAST(ast_pp)
	if *printAst {
		for _, l := range lines {
			fmt.Println(l)
		}
		fmt.Println()
	}
	nodes := convertLinesToNodes(lines)
	tree := buildTree(nodes, 0)

	p := program.NewProgram()
	err = transpiler.TranspileAST(cFilePath, p, tree[0].(ast.Node))
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, p.FileSet, p.File); err != nil {
		panic(err)
	}

	return buf.String()
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file.c>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *version {
		// Simply print out the version and exit.
		fmt.Println(Version)
		return
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Print(Start(flag.Args()))
}
