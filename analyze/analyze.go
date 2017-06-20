package analyze

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/transpiler"
)

// ProgramArgs - argument of program
type ProgramArgs struct {
	Verbose     bool
	Ast         bool
	InputFile   string
	OutputFile  string
	PackageName string
}

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

// ToJSON - tree convert to JSON
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

// Start - base function
func Start(args ProgramArgs) error {
	if os.Getenv("GOPATH") == "" {
		panic("The $GOPATH must be set.")
	}

	// 1. Compile it first (checking for errors)
	_, err := os.Stat(args.InputFile)
	if err != nil {
		return fmt.Errorf("Input file is not found")
	}

	// 2. Preprocess
	var pp []byte
	{
		cmd := exec.Command("clang", "-E", args.InputFile)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("preprocess failed: %v\nStdErr = %v", err, stderr.String())
		}
		pp = []byte(out.String())
	}

	ppFilePath := path.Join(os.TempDir(), "pp.c")
	err = ioutil.WriteFile(ppFilePath, pp, 0644)
	if err != nil {
		return fmt.Errorf("writing to /tmp/pp.c failed: %v", err)
	}

	// 3. Generate JSON from AST
	astPP, err := exec.Command("clang", "-Xclang", "-ast-dump", "-fsyntax-only", ppFilePath).Output()
	if err != nil {
		// If clang fails it still prints out the AST, so we have to run it
		// again to get the real error.
		errBody, _ := exec.Command("clang", ppFilePath).CombinedOutput()

		panic("clang failed: " + err.Error() + ":\n\n" + string(errBody))
	}

	lines := readAST(astPP)
	if args.Ast {
		for _, l := range lines {
			fmt.Println(l)
		}
		fmt.Println()
	}

	nodes := convertLinesToNodes(lines)
	tree := buildTree(nodes, 0)

	p := program.NewProgram()
	p.Verbose = args.Verbose

	err = transpiler.TranspileAST(args.InputFile, args.PackageName, p, tree[0].(ast.Node))
	if err != nil {
		panic(err)
	}

	outputFilePath := args.OutputFile

	if outputFilePath == "" {
		cleanFileName := filepath.Clean(filepath.Base(args.InputFile))
		extension := filepath.Ext(args.InputFile)

		outputFilePath = cleanFileName[0:len(cleanFileName)-len(extension)] + ".go"
	}

	err = ioutil.WriteFile(outputFilePath, []byte(p.String()), 0755)
	if err != nil {
		return fmt.Errorf("writing C output file failed: %v", err)
	}

	return nil
}
