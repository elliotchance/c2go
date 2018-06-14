package darwin

import "github.com/elliotchance/c2go/noarch"

// BuiltinVsprintfChk - implementation __builtin___vsprintf_chk
func BuiltinVsprintfChk(buffer *byte, _ int32, n int32, format *byte, args noarch.VaList) int32 {
	return noarch.Sprintf(buffer, format, args.Args)
}

// BuiltinVsnprintfChk - implementation __builtin___vsnprintf_chk
func BuiltinVsnprintfChk(buffer *byte, n int32, _ int32, _ int32, format *byte, args noarch.VaList) int32 {
	return noarch.Sprintf(buffer, format, args.Args)
}

// BuiltinSprintfChk - implementation __builtin___sprintf_chk
func BuiltinSprintfChk(buffer *byte, _ int32, n int32, format *byte, args ...interface{}) int32 {
	return noarch.Sprintf(buffer, format, args)
}

// BuiltinSnprintfChk - implementation __builtin___snprintf_chk
func BuiltinSnprintfChk(buffer *byte, n int32, _ int32, _ int32, format *byte, args ...interface{}) int32 {
	return noarch.Sprintf(buffer, format, args)
}
