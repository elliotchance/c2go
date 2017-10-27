// +build ignore
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

/*
see https://github.com/elliotchance/c2go/issues/143
Example:

// Define a value for each type.
int i = 10;
float f = 3.2;

// All combinations of `+` would be:
ok_eq(i + i, 20);
ok_eq(i + f, 13.2);
ok_eq(f + i, 13.2);
ok_eq(f + f, 6.4);
*/

func GenerateCombinations(alphabet string, length int) <-chan string {
	c := make(chan string)

	// Starting a separate goroutine that will create all the combinations,
	// feeding them to the channel c
	go func(c chan string) {
		defer close(c) // Once the iteration function is finished, we close the channel

		AddLetter(c, "", alphabet, length) // We start by feeding it an empty string
	}(c)

	return c // Return the channel to the calling function
}

// AddLetter adds a letter to the combination to create a new combination.
// This new combination is passed on to the channel before we call AddLetter once again
// to add yet another letter to the new combination in case length allows it
func AddLetter(c chan string, combo string, alphabet string, length int) {
	// Check if we reached the length limit
	// If so, we just return without adding anything
	if length <= 0 {
		return
	}

	var newCombo string
	for _, ch := range alphabet {
		newCombo = combo + string(ch)
		c <- newCombo
		AddLetter(c, newCombo, alphabet, length-1)
	}
}

type Value struct {
	Type           string
	Name           string
	DefaultValue   string
	PrintSpecifier string
}

func (v Value) String() string {
	return v.Type + " " + v.Name + " = " + v.DefaultValue + ";"
}

type Combination struct {
	values string
}

func main() {
	// All checked type
	Types := []Value{
		// TODO : Don't understood - How to add char ?
		//	Value{"char", "c", "'a'", "%c"},
		Value{"int", "i", "10", "%d"},
		Value{"float", "f", "12.222", "%f"},
		Value{"double", "d", "2.333", "%f"},
	}

	// Checking
	for _, t := range Types {
		if len(t.Name) != 1 {
			panic("Name of types must have only one symbol")
		}
	}

	// Create one string line of types names
	var types string
	for _, t := range Types {
		types += t.Name
	}

	var combinations []string
	for c := range GenerateCombinations(types, len(types)) {
		combinations = append(combinations, c)
	}

	// String line of code
	var code string

	code += `
// DO NOT EDIT - FILE GENERATED
#include <stdio.h>
#include "tests.h"

int main(){
`
	// initialize the types
	for _, t := range Types {
		code += fmt.Sprintf("\t%s\n", t)
	}

	// amount of tests
	code += fmt.Sprintf("	plan(%d);\n", len(combinations)*len(Types))

	for _, specifierType := range Types {
		specifier := specifierType.PrintSpecifier

		code += fmt.Sprintf("	diag(\"%v\");\n", strings.Replace(specifier, "%", "0/0 ", -1))

		for _, combination := range combinations {
			var comb string
			for i, c := range combination {
				if i == 0 {
					comb += fmt.Sprintf(" %c", c)
					continue
				}
				comb += fmt.Sprintf(" + %c", c)
			}
			result := getResultFromC(Types, combination, specifier)
			code += fmt.Sprintf("	is_eq( ( %s ) (%s), %v);\n", specifierType.Type, comb, result)
		}
	}
	code += fmt.Sprintf(`

	done_testing();
}`)
	err := ioutil.WriteFile("./tests/operators_generated.c", []byte(code), 0644)
	if err != nil {
		panic(err)
	}
}

func getResultFromC(v []Value, combination string, specifier string) string {

	name := fmt.Sprintf("gcc%s", combination)

	// generate C code
	var code string
	code += `
			#include <stdio.h>

			int main(){
				`
	// initialize the types
	for _, t := range v {
		code += fmt.Sprintf("\t%s\n", t)
	}

	// create a body
	var comb string
	for i, c := range combination {
		if i == 0 {
			comb += fmt.Sprintf(" %c", c)
			continue
		}
		comb += fmt.Sprintf(" + %c", c)
	}
	code += fmt.Sprintf("	printf(\"%s\",%s);\n", specifier, comb)

	// finalize
	code += `
			return 0;
		}`
	// write the file
	err := ioutil.WriteFile("/tmp/"+name+".c", []byte(code), 0644)
	if err != nil {
		panic(err)
	}

	// execute in gcc
	{
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command("gcc", "-o", "/tmp/"+name, "/tmp/"+name+".c")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			panic(fmt.Errorf("Err = %v\nStdErr = %v", err, stderr.String()))
		}
	}

	// run result
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("/tmp/" + name)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		panic(fmt.Errorf("Err = %v\nStdErr = %v", err, stderr.String()))
	}
	return out.String()
}
