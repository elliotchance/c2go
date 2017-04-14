package main

import (
	"regexp"
	"strings"
)

// FunctionDefinition contains the prototype definition for a function.
type FunctionDefinition struct {
	Name          string
	ReturnType    string
	ArgumentTypes []string
}

var functionDefinitions = map[string]FunctionDefinition{}

var builtInFunctionDefinitionsHaveBeenLoaded = false

var builtInFunctionDefinitions = []string{
	// darwin/assert.h
	"int __builtin_expect(int, int)",
	"bool __assert_rtn(const char*, const char*, int, const char*)",

	// darwin/ctype.h
	"uint32 __istype(__darwin_ct_rune_t, uint32)",
	"__darwin_ct_rune_t __isctype(__darwin_ct_rune_t, uint32)",
	"__darwin_ct_rune_t __tolower(__darwin_ct_rune_t)",
	"__darwin_ct_rune_t __toupper(__darwin_ct_rune_t)",
	"uint32 __maskrune(__darwin_ct_rune_t, uint32)",

	// darwin/math.h
	"double __builtin_fabs(double)",
	"float __builtin_fabsf(float)",
	"double __builtin_fabsl(double)",
	"double __builtin_inf()",
	"float __builtin_inff()",
	"double __builtin_infl()",
	"Double2 __sincospi_stret(double)",
	"Float2 __sincospif_stret(float)",
	"Double2 __sincos_stret(double)",
	"Float2 __sincosf_stret(float)",

	// linux/assert.h
	"bool __assert_fail(const char*, const char*, unsigned int, const char*)",
}

// getFunctionDefinition will return nil if the function does not exist (is not
// registered).
func getFunctionDefinition(functionName string) FunctionDefinition {
	loadFunctionDefinitions()

	return functionDefinitions[functionName]
}

// addFunctionDefinition registers a function definition. If the definition
// already exists it will be replaced.
func addFunctionDefinition(f FunctionDefinition) {
	loadFunctionDefinitions()

	functionDefinitions[f.Name] = f
}

func loadFunctionDefinitions() {
	if builtInFunctionDefinitionsHaveBeenLoaded {
		return
	}

	builtInFunctionDefinitionsHaveBeenLoaded = true

	for _, f := range builtInFunctionDefinitions {
		match := regexp.MustCompile(`^(.+) (.+)\((.*)\)$`).
			FindStringSubmatch(f)

		argumentTypes := strings.Split(match[3], ",")
		for i := range argumentTypes {
			argumentTypes[i] = strings.TrimSpace(argumentTypes[i])
		}

		addFunctionDefinition(FunctionDefinition{
			Name:          match[2],
			ReturnType:    match[1],
			ArgumentTypes: argumentTypes,
		})
	}
}
