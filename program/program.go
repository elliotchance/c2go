// Package program contains high-level orchestration and state of the input and
// output program during transpilation.
package program

import (
	"bytes"
	"fmt"
	"go/format"
	"go/token"

	goast "go/ast"

	"strings"

	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/util"
)

// StructRegistry is a map of Struct for struct types and union type
type StructRegistry map[string]*Struct

// HasType method check if type exists
func (sr StructRegistry) HasType(typename string) bool {
	_, exists := sr[typename]

	return exists
}

// Program contains all of the input, output and transpition state of a C
// program to a Go program.
type Program struct {
	// All of the Go import paths required for this program.
	imports []string

	// These are for the output Go AST.
	FileSet *token.FileSet
	File    *goast.File

	// One a type is defined it will be ignored if a future type of the same
	// name appears.
	typesAlreadyDefined []string

	// Contains the current function name during the transpilation.
	Function *ast.FunctionDecl

	functionDefinitions                      map[string]FunctionDefinition
	builtInFunctionDefinitionsHaveBeenLoaded bool

	// These are used to setup the runtime before the application begins. An
	// example would be to setup globals with stdin file pointers on certain
	// platforms.
	startupStatements []goast.Stmt

	// This is used to generate globally unique names for temporary variables
	// and other generated code. See GetNextIdentifier().
	nextUniqueIdentifier int

	// The definitions for defined structs.
	// TODO: This field should be protected through proper getters and setters.
	Structs StructRegistry
	Unions  StructRegistry

	// If verbose is on progress messages will be printed immediately as code
	// comments (so that they do not interfere with the program output).
	Verbose bool

	// Contains the messages (for example, "// Warning") generated when
	// transpiling the AST. These messages, which are code comments, are
	// appended to the very top of the output file. See AddMessage().
	messages []string

	// messagePosition - position of slice messages, added like a comment
	// in output Go code
	messagePosition int

	// A map of all the global variables (variables that exist outside of a
	// function) and their types.
	GlobalVariables map[string]string

	// This option is not available through the command line. It is to allow the
	// internal integration testing to generate the output in the form of a
	// Go-test rather than a standalone Go file.
	OutputAsTest bool

	// EnumConstantToEnum - a map with key="EnumConstant" and value="enum type"
	// clang don`t show enum constant with enum type,
	// so we have to use hack for repair the type
	EnumConstantToEnum map[string]string

	// EnumTypedefName - a map with key="Name of typedef enum" and
	// value="exist ot not"
	EnumTypedefName map[string]bool

	// TypedefType - map for type alias, for example:
	// C  : typedef int INT;
	// Map: key = INT, value = int
	// Important: key and value are C types
	TypedefType map[string]string

	// Comments
	Comments []Comment

	// commentLine - a map with:
	// key    - filename
	// value  - last comment inserted in Go code
	commentLine map[string]int

	// IncludeHeaders - list of C header
	IncludeHeaders []IncludeHeader

	// NodeMap - a map containing all the program's nodes with:
	// key    - the node address
	// value  - the node
	NodeMap map[ast.Address]ast.Node
}

// Comment - position of line comment '//...'
type Comment struct {
	File    string
	Line    int
	Comment string
}

// IncludeHeader - struct for C include header
type IncludeHeader struct {
	HeaderName   string
	IsUserSource bool
}

// NewProgram creates a new blank program.
func NewProgram() (p *Program) {
	defer func() {
		// Need for "stdbool.h"
		p.TypedefType["_Bool"] = "signed char"
	}()
	return &Program{
		imports:             []string{},
		typesAlreadyDefined: []string{},
		startupStatements:   []goast.Stmt{},
		Structs: StructRegistry(map[string]*Struct{
			// Structs without implementations inside system C headers
			// Example node for adding:
			// &ast.TypedefDecl{ ... Type:"struct __locale_struct *" ... }

			"struct __va_list_tag [1]": {
				Name:    "struct __va_list_tag [1]",
				IsUnion: false,
			},

			// Pos:ast.Position{File:"/usr/include/xlocale.h", Line:27
			"struct __locale_struct *": {
				Name:    "struct __locale_struct *",
				IsUnion: false,
			},

			// Pos:ast.Position{File:"/usr/include/x86_64-linux-gnu/sys/time.h", Line:61
			"struct timezone *__restrict": {
				Name:    "struct timezone *__restrict",
				IsUnion: false,
			},
		}),
		Unions:              make(StructRegistry),
		Verbose:             false,
		messages:            []string{},
		GlobalVariables:     map[string]string{},
		EnumConstantToEnum:  map[string]string{},
		EnumTypedefName:     map[string]bool{},
		TypedefType:         map[string]string{},
		commentLine:         map[string]int{},
		IncludeHeaders:      []IncludeHeader{},
		functionDefinitions: map[string]FunctionDefinition{},
		NodeMap:             map[ast.Address]ast.Node{},
		builtInFunctionDefinitionsHaveBeenLoaded: false,
	}
}

// AddMessage adds a message (such as a warning or error) comment to the output
// file. Usually the message is generated from one of the Generate functions in
// the ast package.
//
// It is expected that the message already have the comment ("//") prefix.
//
// The message will not be appended if it is blank. This is because the Generate
// functions return a blank string conditionally when there is no error.
//
// The return value will be true if a message was added, otherwise false.
func (p *Program) AddMessage(message string) bool {
	if message == "" {
		return false
	}

	p.messages = append(p.messages, message)

	// Compactizarion warnings stack
	if len(p.messages) > 1 {
		var (
			new  = len(p.messages) - 1
			last = len(p.messages) - 2
		)
		// Warning collapsing for minimaze warnings
		warning := "// Warning"
		if strings.HasPrefix(p.messages[last], warning) {
			l := p.messages[last][len(warning):]
			if strings.HasSuffix(p.messages[new], l) {
				p.messages[last] = p.messages[new]
				p.messages = p.messages[0:new]
			}
		}
	}

	return true
}

// GetMessageComments - get messages "Warnings", "Error" like a comment
// Location of comments only NEAR of error or warning and
// don't show directly location
func (p *Program) GetMessageComments() (_ *goast.CommentGroup) {
	var group goast.CommentGroup
	if p.messagePosition < len(p.messages) {
		for i := p.messagePosition; i < len(p.messages); i++ {
			group.List = append(group.List, &goast.Comment{
				Text: p.messages[i],
			})
		}
		p.messagePosition = len(p.messages)
	}
	return &group
}

// GetComments - return comments
func (p *Program) GetComments(n ast.Position) (out []*goast.Comment) {
	beginLine := p.commentLine[n.File]
	lastLine := n.LineEnd
	for i := range p.Comments {
		if p.Comments[i].File == n.File {
			if beginLine < p.Comments[i].Line && p.Comments[i].Line <= lastLine {
				out = append(out, &goast.Comment{
					Text: p.Comments[i].Comment,
				})
				if p.Comments[i].Comment[0:2] == "/*" {
					out = append(out, &goast.Comment{
						Text: "// ",
					})
				}
			}
		}
	}
	if len(out) > 0 {
		out = append(out, &goast.Comment{
			Text: "// ",
		})
	}
	p.commentLine[n.File] = lastLine
	return
}

// GetStruct returns a struct object (representing struct type or union type) or
// nil if doesn't exist. This method can get struct or union in the same way and
// distinguish only by the IsUnion field. `name` argument is the C like
// `struct a_struct`, it allow pointer type like `union a_union *`. Pointer
// types used in a DeclRefExpr in the case a deferenced structure by using `->`
// operator to access to a field like this: a_struct->member .
//
// This method is used in collaboration with the field
// "c2go/program".*Struct.IsUnion to simplify the code like in function
// "c2go/transpiler".transpileMemberExpr() where the same *Struct value returned
// by this method is used in the 2 cases, in the case where the value has a
// struct type and in the case where the value has an union type.
func (p *Program) GetStruct(name string) *Struct {
	if name == "" {
		return nil
	}

	last := len(name) - 1

	// That allow to get struct from pointer type
	if name[last] == '*' {
		name = name[:last]
	}

	name = strings.TrimSpace(name)

	res, ok := p.Structs[name]
	if ok {
		return res
	}

	return p.Unions[name]
}

// IsTypeAlreadyDefined will return true if the typeName has already been
// defined.
//
// A type could be defined:
//
// 1. Initially. That is, before the transpilation starts (hard-coded).
// 2. By calling DefineType throughout the transpilation.
func (p *Program) IsTypeAlreadyDefined(typeName string) bool {
	return util.InStrings(typeName, p.typesAlreadyDefined)
}

// DefineType will record a type as having already been defined. The purpose for
// this is to not generate Go for a type more than once. C allows variables and
// other entities (such as function prototypes) to be defined more than once in
// some cases. An example of this would be static variables or functions.
func (p *Program) DefineType(typeName string) {
	p.typesAlreadyDefined = append(p.typesAlreadyDefined, typeName)
}

// GetNextIdentifier generates a new globally unique identifier name. This can
// be used for variables and functions in generated code.
//
// The value of prefix is only useful for readability in the code. If the prefix
// is an empty string then the prefix "__temp" will be used.
func (p *Program) GetNextIdentifier(prefix string) string {
	if prefix == "" {
		prefix = "temp"
	}

	identifierName := fmt.Sprintf("%s%d", prefix, p.nextUniqueIdentifier)
	p.nextUniqueIdentifier++

	return identifierName
}

// String generates the whole output Go file as a string. This will include the
// messages at the top of the file and all the rendered Go code.
func (p *Program) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(`/*
	Package %s - transpiled by c2go version: %s

	If you have found any issues, please raise an issue at:
	https://github.com/elliotchance/c2go/
*/

`, p.File.Name.Name, Version))

	// First write all the messages. The double newline afterwards is important
	// so that the package statement has a newline above it so that the warnings
	// are not part of the documentation for the package.
	buf.WriteString(strings.Join(p.messages, "\n") + "\n\n")

	if err := format.Node(&buf, p.FileSet, p.File); err != nil {
		// Printing the entire AST will generate a lot of output. However, it is
		// the only way to debug this type of error. Hopefully the error
		// (printed immediately afterwards) will give a clue.
		//
		// You may see an error like:
		//
		//     panic: format.Node internal error (692:23: expected selector or
		//     type assertion, found '[')
		//
		// This means that when Go was trying to convert the Go AST to source
		// code it has come across a value or attribute that is illegal.
		//
		// The line number it is referring to (in this case, 692) is not helpful
		// as it references the internal line number of the Go code which you
		// will never see.
		//
		// The "[" means that there is a bracket in the wrong place. Almost
		// certainly in an identifer, like:
		//
		//     noarch.IntTo[]byte("foo")
		//
		// The "[]" which is obviously not supposed to be in the function name
		// is causing the syntax error. However, finding the original code that
		// produced this can be tricky.
		//
		// The first step is to filter down the AST output to probably lines.
		// In the error message it said that there was a misplaced "[" so that's
		// what we will search for. Using the original command (that generated
		// thousands of lines) we will add two grep filters:
		//
		//     go test ... | grep "\[" | grep -v '{$'
		//     #                   |     |
		//     #                   |     ^ This excludes lines that end with "{"
		//     #                   |       which almost certainly won't be what
		//     #                   |       we are looking for.
		//     #                   |
		//     #                   ^ This is the character we are looking for.
		//
		// Hopefully in the output you should see some lines, like (some lines
		// removed for brevity):
		//
		//     9083  .  .  .  .  .  .  .  .  .  .  Name: "noarch.[]byteTo[]int"
		//     9190  .  .  .  .  .  .  .  .  .  Name: "noarch.[]intTo[]byte"
		//
		// These two lines are clearly the error because a name should not look
		// like this.
		//
		// Looking at the full output of the AST (thousands of lines) and
		// looking at those line numbers should give you a good idea where the
		// error is coming from; by looking at the parents of the bad lines.
		_ = goast.Print(p.FileSet, p.File)

		panic(err)
	}

	// Add comments at the end C file
	for file, beginLine := range p.commentLine {
		for i := range p.Comments {
			if p.Comments[i].File == file {
				if beginLine < p.Comments[i].Line {
					buf.WriteString(fmt.Sprintln(p.Comments[i].Comment))
				}
			}
		}
	}

	// simplify Go code. Example :
	// Before:
	// func compare(a interface {
	// }, b interface {
	// }) (c2goDefaultReturn int) {
	// After :
	// func compare(a interface {}, b interface {}) (c2goDefaultReturn int) {
	reg := util.GetRegex("interface( )?{(\r*)\n(\t*)}")

	return string(reg.ReplaceAll(buf.Bytes(), []byte("interface {}")))
}

// IncludeHeaderIsExists - return true if C #include header is inside list
func (p *Program) IncludeHeaderIsExists(includeHeader string) bool {
	for _, inc := range p.IncludeHeaders {
		if strings.HasSuffix(inc.HeaderName, includeHeader) {
			return true
		}
	}
	return false
}

// SetNodes will add the given nodes and all their children to the program's node map.
func (p *Program) SetNodes(nodes []ast.Node) {
	for _, n := range nodes {
		if n == nil {
			continue
		}
		var setNode = true
		addr := n.Address()
		if addr == 0 {
			setNode = false
		} else if _, ok := n.(*ast.Record); ok {
			setNode = false
		}
		if setNode {
			p.NodeMap[addr] = n
		}
		p.SetNodes(n.Children())
	}
}

// DeclareType defines a type without adding the real definition of the type to the list
// of declarations. This is useful for types which are defined in a header file, but are
// implemented in another *.c file.
// This way they can already be used for variable declarations and sizeof.
func (p *Program) DeclareType(n *ast.RecordDecl, correctName string) (err error) {
	name := correctName

	// TODO: Some platform structs are ignored.
	// https://github.com/elliotchance/c2go/issues/85
	if name == "__locale_struct" ||
		name == "__sigaction" ||
		name == "sigaction" {
		err = nil
		return
	}

	s := NewStruct(n)
	if s.IsUnion {
		p.Unions["union "+name] = s
	} else {
		p.Structs["struct "+name] = s
	}

	return
}
