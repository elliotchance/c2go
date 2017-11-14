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
}

// NewProgram creates a new blank program.
func NewProgram() *Program {
	return &Program{
		imports:             []string{},
		typesAlreadyDefined: []string{},
		startupStatements:   []goast.Stmt{},
		Structs: StructRegistry(map[string]*Struct{
			// Structs without implementations inside system C headers
			// Example node for adding:
			// &ast.TypedefDecl{ ... Type:"struct __locale_struct *" ... }

			"struct __va_list_tag [1]": &Struct{
				Name:    "struct __va_list_tag [1]",
				IsUnion: false,
			},

			// Pos:ast.Position{File:"/usr/include/xlocale.h", Line:27
			"struct __locale_struct *": &Struct{
				Name:    "struct __locale_struct *",
				IsUnion: false,
			},

			// Pos:ast.Position{File:"/usr/include/x86_64-linux-gnu/sys/time.h", Line:61
			"struct timezone *__restrict": &Struct{
				Name:    "struct timezone *__restrict",
				IsUnion: false,
			},
		}),
		Unions:             make(StructRegistry),
		Verbose:            false,
		messages:           []string{},
		GlobalVariables:    map[string]string{},
		EnumConstantToEnum: map[string]string{},
		EnumTypedefName:    map[string]bool{},
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
	return true
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
		goast.Print(p.FileSet, p.File)

		panic(err)
	}

	return buf.String()
}
