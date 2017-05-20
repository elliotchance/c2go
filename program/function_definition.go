package program

import (
	"regexp"
	"strings"
)

// FunctionDefinition contains the prototype definition for a function.
type FunctionDefinition struct {
	// The name of the function, like "printf".
	Name string

	// The C return type, like "int".
	ReturnType string

	// The C argument types, like ["bool", "int"]. There is currently no way
	// to represent a varargs.
	ArgumentTypes []string

	// If this is not empty then this function name should be used instead
	// of the Name. Many low level functions have an exact match with a Go
	// function. For example, "sin()".
	Substitution string

	// Can be overriden with the substitution to rearrange the return variables
	// and parameters. When either of these are nil the behavior is to keep the
	// single return value and parameters the same.
	ReturnParameters []int
	Parameters       []int
}

var functionDefinitions map[string]FunctionDefinition

var builtInFunctionDefinitionsHaveBeenLoaded = false

// Each of the predefined function have a syntax that allows them to be easy to
// read (and maintain). For example:
//
//     double __builtin_fabs(double) -> darwin.Fabs
//
// Declares the prototype of __builtin_fabs (a low level function implemented
// only on Mac) with a specific substitution provided. This means that it should
// replace any instance of __builtin_fabs with:
//
//     github.com/elliotchance/c2go/darwin.Fabs
//
// The substitution is optional.
//
// The substituted function can also move the parameters and return value
// positions. This is called a transformation. For example:
//
//     size_t fread(void*, size_t, size_t, FILE*) -> $0, $1 = noarch.Fread($2, $3, $4)
//
// Where $0 represents the C return value and $1 and above are for each of the
// parameters.
//
// Transformations can also be used to specify varaible that need to be passed
// by reference by using the prefix "&" instead of "$":
//
//     size_t fread(void*, size_t, size_t, FILE*) -> $0 = noarch.Fread(&1, $2, $3, $4)
//
var builtInFunctionDefinitions = []string{
	// darwin/assert.h
	"int __builtin_expect(int, int) -> darwin.BuiltinExpect",
	"bool __assert_rtn(const char*, const char*, int, const char*) -> darwin.AssertRtn",

	// darwin/ctype.h
	"uint32 __istype(__darwin_ct_rune_t, uint32) -> darwin.IsType",
	"__darwin_ct_rune_t __isctype(__darwin_ct_rune_t, uint32) -> darwin.IsCType",
	"__darwin_ct_rune_t __tolower(__darwin_ct_rune_t) -> darwin.ToLower",
	"__darwin_ct_rune_t __toupper(__darwin_ct_rune_t) -> darwin.ToUpper",
	"uint32 __maskrune(__darwin_ct_rune_t, uint32) -> darwin.MaskRune",

	// linux/ctype.h
	"const unsigned short int** __ctype_b_loc() -> linux.CtypeLoc",
	"int tolower(int) -> linux.ToLower",
	"int toupper(int) -> linux.ToUpper",

	// darwin/math.h
	"double __builtin_fabs(double) -> darwin.Fabs",
	"float __builtin_fabsf(float) -> darwin.Fabsf",
	"double __builtin_fabsl(double) -> darwin.Fabsl",
	"double __builtin_inf() -> darwin.Inf",
	"float __builtin_inff() -> darwin.Inff",
	"double __builtin_infl() -> darwin.Infl",
	"Double2 __sincospi_stret(double) -> darwin.SincospiStret",
	"Float2 __sincospif_stret(float) -> darwin.SincospifStret",
	"Double2 __sincos_stret(double) -> darwin.SincosStret",
	"Float2 __sincosf_stret(float) -> darwin.SincosfStret",

	// linux/assert.h
	"bool __assert_fail(const char*, const char*, unsigned int, const char*) -> linux.AssertFail",

	// linux/stdio.h
	"int _IO_getc(FILE*) -> noarch.Fgetc",
	"int _IO_putc(int, FILE*) -> noarch.Fputc",

	// math.h
	"double acos(double) -> math.Acos",
	"double asin(double) -> math.Asin",
	"double atan(double) -> math.Atan",
	"double atan2(double) -> math.Atan2",
	"double ceil(double) -> math.Ceil",
	"double cos(double) -> math.Cos",
	"double cosh(double) -> math.Cosh",
	"double exp(double) -> math.Exp",
	"double fabs(double) -> math.Abs",
	"double floor(double) -> math.Floor",
	"double fmod(double) -> math.Mod",
	"double ldexp(double) -> math.Ldexp",
	"double log(double) -> math.Log",
	"double log10(double) -> math.Log10",
	"double pow(double) -> math.Pow",
	"double sin(double) -> math.Sin",
	"double sinh(double) -> math.Sinh",
	"double sqrt(double) -> math.Sqrt",
	"double tan(double) -> math.Tan",
	"double tanh(double) -> math.Tanh",

	// stdio.h
	"int printf(const char*) -> noarch.Printf",
	"int scanf(const char*) -> noarch.Scanf",
	"int putchar(int) -> darwin.Putchar",
	"int puts(const char *) -> noarch.Puts",
	"FILE* fopen(const char *, const char *) -> noarch.Fopen",
	"int fclose(FILE*) -> noarch.Fclose",
	"int remove(const char*) -> noarch.Remove",
	"int rename(const char*, const char*) -> noarch.Rename",
	"int fputs(const char*, FILE*) -> noarch.Fputs",
	"FILE* tmpfile() -> noarch.Tmpfile",
	"char* fgets(char*, int, FILE*) -> noarch.Fgets",
	"void rewind(FILE*) -> noarch.Rewind",
	"int feof(FILE*) -> noarch.Feof",
	"char* tmpnam(char*) -> noarch.Tmpnam",
	"int fflush(FILE*) -> noarch.Fflush",
	"int fprintf(FILE*, const char*) -> noarch.Fprintf",
	"int fscanf(FILE*, const char*) -> noarch.Fscanf",
	"int fgetc(FILE*) -> noarch.Fgetc",
	"int fputc(int, FILE*) -> noarch.Fputc",
	"int getc(FILE*) -> noarch.Fgetc",
	"int getchar() -> noarch.Getchar",
	"int putc(int, FILE*) -> noarch.Fputc",
	"int fseek(FILE*, long int, int) -> noarch.Fseek",
	"long ftell(FILE*) -> noarch.Ftell",
	"int fread(void*, int, int, FILE*) -> $0 = noarch.Fread(&1, $2, $3, $4)",
	"int fwrite(char*, int, int, FILE*) -> noarch.Fwrite",
	"int fgetpos(FILE*, int*) -> noarch.Fgetpos",
	"int fsetpos(FILE*, int*) -> noarch.Fsetpos",

	// string.h
	"size_t strlen(const char*) -> noarch.Strlen",

	// stdlib.h
	"int atoi(const char*) -> noarch.Atoi",
	"long strtol(const char *, char **, int) -> noarch.Strtol",
	"void* malloc(unsigned int) -> noarch.Malloc",
	"void free(void*) -> noarch.Free",

	// I'm not sure which header file these comes from?
	"uint32 __builtin_bswap32(uint32) -> darwin.BSwap32",
	"uint64 __builtin_bswap64(uint64) -> darwin.BSwap64",
}

// GetFunctionDefinition will return nil if the function does not exist (is not
// registered).
func GetFunctionDefinition(functionName string) *FunctionDefinition {
	loadFunctionDefinitions()

	if f, ok := functionDefinitions[functionName]; ok {
		return &f
	}

	return nil
}

// AddFunctionDefinition registers a function definition. If the definition
// already exists it will be replaced.
func AddFunctionDefinition(f FunctionDefinition) {
	loadFunctionDefinitions()

	functionDefinitions[f.Name] = f
}

// dollarArgumentsToIntSlice converts a list of dollar arguments, like "$1, &2"
// into a slice of integers; [1, -2].
//
// This function requires at least one argument in s, but only arguments upto
// $9 or &9.
func dollarArgumentsToIntSlice(s string) []int {
	r := []int{}
	multiplier := 1

	for _, c := range s {
		if c == '$' {
			multiplier = 1
		}
		if c == '&' {
			multiplier = -1
		}

		if c >= '0' && c <= '9' {
			r = append(r, multiplier*(int(c)-'0'))
		}
	}

	return r
}

func loadFunctionDefinitions() {
	if builtInFunctionDefinitionsHaveBeenLoaded {
		return
	}

	functionDefinitions = map[string]FunctionDefinition{}
	builtInFunctionDefinitionsHaveBeenLoaded = true

	for _, f := range builtInFunctionDefinitions {
		match := regexp.MustCompile(`^(.+) ([^ ]+)\(([, a-z*A-Z_0-9]*)\)( -> .+)?$`).
			FindStringSubmatch(f)

		// Unpack argument types.
		argumentTypes := strings.Split(match[3], ",")
		for i := range argumentTypes {
			argumentTypes[i] = strings.TrimSpace(argumentTypes[i])
		}
		if len(argumentTypes) == 1 && argumentTypes[0] == "" {
			argumentTypes = []string{}
		}

		// Defaults for transformations.
		var returnParameters, parameters []int

		// Substitution rules.
		substitution := match[4]
		if substitution != "" {
			substitution = strings.TrimLeft(substitution, " ->")

			// The substitution might also rearrange the parameters (return and
			// parameter transformation).
			subMatch := regexp.MustCompile(`^(.*?) = (.*)\((.*)\)$`).
				FindStringSubmatch(substitution)
			if len(subMatch) > 0 {
				returnParameters = dollarArgumentsToIntSlice(subMatch[1])
				parameters = dollarArgumentsToIntSlice(subMatch[3])
				substitution = subMatch[2]
			}
		}

		if strings.HasPrefix(substitution, "darwin.") ||
			strings.HasPrefix(substitution, "linux.") ||
			strings.HasPrefix(substitution, "noarch.") {
			substitution = "github.com/elliotchance/c2go/" + substitution
		}

		AddFunctionDefinition(FunctionDefinition{
			Name:             match[2],
			ReturnType:       match[1],
			ArgumentTypes:    argumentTypes,
			Substitution:     substitution,
			ReturnParameters: returnParameters,
			Parameters:       parameters,
		})
	}
}
