package darwin

import "github.com/elliotchance/c2go/noarch"

// BuiltinVsprintfChk - implementation __builtin___vsprintf_chk
func BuiltinVsprintfChk(buffer []byte, _ int, n int, format []byte, args ...interface{}) int {
	return noarch.Sprintf(buffer, format, args)
}

// BuiltinVsnprintfChk - implementation __builtin___vsnprintf_chk
func BuiltinVsnprintfChk(buffer []byte, n int, _ int, _ int, format []byte, args ...interface{}) int {
	return noarch.Sprintf(buffer, format, args)
}

// BuiltinSprintfChk - implementation __builtin___sprintf_chk
func BuiltinSprintfChk(buffer []byte, _ int, n int, format []byte, args ...interface{}) int {
	return noarch.Sprintf(buffer, format, args)
}

// BuiltinSnprintfChk - implementation __builtin___snprintf_chk
func BuiltinSnprintfChk(buffer []byte, n int, _ int, _ int, format []byte, args ...interface{}) int {
	return noarch.Sprintf(buffer, format, args)
}
