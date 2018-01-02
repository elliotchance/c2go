package darwin

import "github.com/elliotchance/c2go/noarch"

// BuiltinSprintfChk - implementation __builtin___sprintf_chk
func BuiltinSprintfChk(buffer []byte, _ int, n int, format []byte, args ...interface{}) int {
	return noarch.Snprintf(buffer, n, format, args)
}
