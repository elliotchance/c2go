package program

import (
	"strings"

	"github.com/elliotchance/c2go/util"
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

	// Can be overridden with the substitution to rearrange the return variables
	// and parameters. When either of these are nil the behavior is to keep the
	// single return value and parameters the same.
	ReturnParameters []int
	Parameters       []int
}

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
// Transformations can also be used to specify variable that need to be passed
// by reference by using the prefix "&" instead of "$":
//
//     size_t fread(void*, size_t, size_t, FILE*) -> $0 = noarch.Fread(&1, $2, $3, $4)
//
var builtInFunctionDefinitions = map[string][]string{
	"assert.h": []string{
		// darwin/assert.h
		"int __builtin_expect(int, int) -> darwin.BuiltinExpect",
		"bool __assert_rtn(const char*, const char*, int, const char*) -> darwin.AssertRtn",

		// linux/assert.h
		"bool __assert_fail(const char*, const char*, unsigned int, const char*) -> linux.AssertFail",
	},
	"ctype.h": []string{
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
	},
	"math.h": []string{
		// linux/math.h
		"int __signbitf(float) -> noarch.Signbitf",
		"int __signbit(double) -> noarch.Signbitd",
		"int __signbitl(long double) -> noarch.Signbitl",
		"int __builtin_signbitf(float) -> noarch.Signbitf",
		"int __builtin_signbit(double) -> noarch.Signbitd",
		"int __builtin_signbitl(long double) -> noarch.Signbitl",
		"int __isnanf(float) -> linux.IsNanf",
		"int __isnan(double) -> noarch.IsNaN",
		"int __isnanl(long double) -> noarch.IsNaN",
		"int __isinff(float) -> linux.IsInff",
		"int __isinf(double) -> linux.IsInf",
		"int __isinfl(long double) -> linux.IsInf",

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
		"float __builtin_huge_valf() -> darwin.Inff",
		"int __inline_signbitf(float) -> noarch.Signbitf",
		"int __inline_signbitd(double) -> noarch.Signbitd",
		"int __inline_signbitl(long double) -> noarch.Signbitl",
		"double __builtin_nanf(const char*) -> darwin.NaN",

		// math.h
		"double acos(double) -> math.Acos",
		"double asin(double) -> math.Asin",
		"double atan(double) -> math.Atan",
		"double atan2(double, double) -> math.Atan2",
		"double ceil(double) -> math.Ceil",
		"double cos(double) -> math.Cos",
		"double cosh(double) -> math.Cosh",
		"double exp(double) -> math.Exp",
		"double fabs(double) -> math.Abs",
		"double floor(double) -> math.Floor",
		"double fmod(double, double) -> math.Mod",
		"double ldexp(double, int) -> noarch.Ldexp",
		"double log(double) -> math.Log",
		"double log10(double) -> math.Log10",
		"double pow(double, double) -> math.Pow",
		"double sin(double) -> math.Sin",
		"double sinh(double) -> math.Sinh",
		"double sqrt(double) -> math.Sqrt",
		"double tan(double) -> math.Tan",
		"double tanh(double) -> math.Tanh",
	},
	"stdio.h": []string{

		// linux/stdio.h
		"int _IO_getc(FILE*) -> noarch.Fgetc",
		"int _IO_putc(int, FILE*) -> noarch.Fputc",

		// stdio.h
		"int printf(const char*) -> noarch.Printf",
		"int scanf(const char*) -> noarch.Scanf",
		"int putchar(int) -> noarch.Putchar",
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
		"int ferror(FILE*) -> noarch.Ferror",
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
		"int fread(void*, int, int, FILE*) -> noarch.Fread",
		"int fwrite(char*, int, int, FILE*) -> noarch.Fwrite",
		"int fgetpos(FILE*, int*) -> noarch.Fgetpos",
		"int fsetpos(FILE*, int*) -> noarch.Fsetpos",
		"int sprintf(char*, const char *) -> noarch.Sprintf",
		"int snprintf(char*, int, const char *) -> noarch.Snprintf",
		"int vsprintf(char*, const char *, struct __va_list_tag *) -> noarch.Vsprintf",
		"int vsnprintf(char*, int, const char *, struct __va_list_tag *) -> noarch.Vsnprintf",
		"void perror(char*) -> noarch.Perror",
		"void clearerr(FILE*) -> noarch.Clearerr",

		// darwin/stdio.h
		"int __builtin___sprintf_chk(char*, int, int, char*) -> darwin.BuiltinSprintfChk",
		"int __builtin___snprintf_chk(char*, int, int, int, char*) -> darwin.BuiltinSnprintfChk",
		"int __builtin___vsprintf_chk(char*, int, int, char *, struct __va_list_tag *) -> darwin.BuiltinVsprintfChk",
		"int __builtin___vsnprintf_chk(char*, int, int, int, char*, struct __va_list_tag *) -> darwin.BuiltinVsnprintfChk",
	},
	"string.h": []string{
		// string.h
		"char* strcasestr(const char*, const char*) -> noarch.Strcasestr",
		"char* strcat(char *, const char *) -> noarch.Strcat",
		"int strcmp(const char *, const char *) -> noarch.Strcmp",
		"char* strerror(int) -> noarch.Strerror",

		// should be: "int strncmp(const char*, const char*, size_t) -> noarch.Strncmp",
		"int strncmp(const char *, const char *, int) -> noarch.Strncmp",
		"char * strchr(char *, int) -> noarch.Strchr",

		"char* strcpy(const char*, char*) -> noarch.Strcpy",
		// should be: "char* strncpy(const char*, char*, size_t) -> noarch.Strncpy",
		"char* strncpy(const char*, char*, int) -> noarch.Strncpy",

		// real return type is "size_t", but it is changed to "int"
		// in according to noarch.Strlen
		"int strlen(const char*) -> noarch.Strlen",

		"char* strstr(const char*, const char*) -> noarch.Strstr",

		// should be: "void* memset(void *, int, size_t) -> noarch.Memset"
		"void* memset(void *, int, int) -> noarch.Memset",

		// should be: "void* memcpy(void *, void *, size_t) -> noarch.Memcpy"
		"void* memcpy(void *, void *, int) -> noarch.Memcpy",

		// should be: "void* memmove(void *, void *, size_t) -> noarch.Memcpy"
		"void* memmove(void *, void *, int) -> noarch.Memcpy",

		// should be: "int memmove(const void *, const void *, size_t) -> noarch.Memcmp"
		"int memcmp(void *, void *, int) -> noarch.Memcmp",

		// darwin/string.h
		// should be: const char*, char*, size_t
		"char* __builtin___strcpy_chk(const char*, char*, int) -> darwin.BuiltinStrcpy",
		// should be: const char*, char*, size_t, size_t
		"char* __builtin___strncpy_chk(const char*, char*, int, int) -> darwin.BuiltinStrncpy",

		// should be: size_t __builtin_object_size(const void*, int)
		"int __builtin_object_size(const char*, int) -> darwin.BuiltinObjectSize",

		// see https://opensource.apple.com/source/Libc/Libc-763.12/include/secure/_string.h.auto.html
		"char* __builtin___strcat_chk(char *, const char *, int) -> darwin.BuiltinStrcat",
		"char* __inline_strcat_chk(char *, const char *) -> noarch.Strcat",
		"void* __builtin___memset_chk(void *, int, int, int) -> darwin.Memset",
		"void* __inline_memset_chk(void *, int, int) -> noarch.Memset",
		"void* __builtin___memcpy_chk(void *, void *, int, int) -> darwin.Memcpy",
		"void* __inline_memcpy_chk(void *, void *, int) -> noarch.Memcpy",
		"void* __builtin___memmove_chk(void *, void *, int, int) -> darwin.Memcpy",
		"void* __inline_memmove_chk(void *, void *, int) -> noarch.Memcpy",
	},
	"stdlib.h": []string{
		// stdlib.h
		"int abs(int) -> noarch.Abs",
		"double atof(const char *) -> noarch.Atof",
		"int atoi(const char*) -> noarch.Atoi",
		"long int atol(const char*) -> noarch.Atol",
		"long long int atoll(const char*) -> noarch.Atoll",
		"div_t div(int, int) -> noarch.Div",
		"void exit(int) -> noarch.Exit",
		"void free(void*) -> noarch.Free",
		"char* getenv(const char *) -> noarch.Getenv",
		"long int labs(long int) -> noarch.Labs",
		"ldiv_t ldiv(long int, long int) -> noarch.Ldiv",
		"long long int llabs(long long int) -> noarch.Llabs",
		"lldiv_t lldiv(long long int, long long int) -> noarch.Lldiv",
		"int rand() -> noarch.Rand",
		// The real definition is srand(unsigned int) however the type would be
		// different. It's easier to change the definition than create a proxy
		// function in stdlib.go.
		"void srand(long long) -> math/rand.Seed",
		"double strtod(const char *, char **) -> noarch.Strtod",
		"float strtof(const char *, char **) -> noarch.Strtof",
		"long strtol(const char *, char **, int) -> noarch.Strtol",
		"long double strtold(const char *, char **) -> noarch.Strtold",
		"long long strtoll(const char *, char **, int) -> noarch.Strtoll",
		"long unsigned int strtoul(const char *, char **, int) -> noarch.Strtoul",
		"long long unsigned int strtoull(const char *, char **, int) -> noarch.Strtoull",
		"void free(void*) -> noarch.Free",
	},
	"syslog.h": []string{
		"void openlog(const char *, int, int) -> noarch.Openlog",
		"int setlogmask(int) -> noarch.Setlogmask",
		"void syslog(int, const char *) -> noarch.Syslog",
		"void vsyslog(int, const char *, struct __va_list_tag *) -> noarch.Vsyslog",
		"void closelog(void) -> noarch.Closelog",
	},
	"time.h": []string{
		// time.h
		"time_t time(time_t *) -> noarch.Time",
		"char* ctime(const time_t *) -> noarch.Ctime",
		"struct tm * localtime(const time_t *) -> noarch.LocalTime",
		"struct tm * gmtime(const time_t *) -> noarch.Gmtime",
		"time_t mktime(struct tm *) -> noarch.Mktime",
		"char * asctime(struct tm *) -> noarch.Asctime",
	},
	"endian.h": []string{
		// I'm not sure which header file these comes from?
		"uint32 __builtin_bswap32(uint32) -> darwin.BSwap32",
		"uint64 __builtin_bswap64(uint64) -> darwin.BSwap64",
	},
	"errno.h": []string{
		// linux
		"int * __errno_location() -> noarch.Errno",

		// darwin
		"int * __error() -> noarch.Errno",
	},
}

// GetFunctionDefinition will return nil if the function does not exist (is not
// registered).
func (p *Program) GetFunctionDefinition(functionName string) *FunctionDefinition {
	p.loadFunctionDefinitions()

	if f, ok := p.functionDefinitions[functionName]; ok {
		return &f
	}

	return nil
}

// AddFunctionDefinition registers a function definition. If the definition
// already exists it will be replaced.
func (p *Program) AddFunctionDefinition(f FunctionDefinition) {
	p.loadFunctionDefinitions()

	p.functionDefinitions[f.Name] = f
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

func (p *Program) loadFunctionDefinitions() {
	if p.builtInFunctionDefinitionsHaveBeenLoaded {
		return
	}

	p.functionDefinitions = map[string]FunctionDefinition{}
	p.builtInFunctionDefinitionsHaveBeenLoaded = true

	for k, v := range builtInFunctionDefinitions {
		if !p.IncludeHeaderIsExists(k) {
			continue
		}

		for _, f := range v {
			match := util.GetRegex(`^(.+) ([^ ]+)\(([, a-z*A-Z_0-9]*)\)( -> .+)?$`).
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
				subMatch := util.GetRegex(`^(.*?) = (.*)\((.*)\)$`).
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

			p.AddFunctionDefinition(FunctionDefinition{
				Name:             match[2],
				ReturnType:       match[1],
				ArgumentTypes:    argumentTypes,
				Substitution:     substitution,
				ReturnParameters: returnParameters,
				Parameters:       parameters,
			})
		}
	}
}
