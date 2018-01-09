package darwin

import "github.com/elliotchance/c2go/noarch"

// BuiltinStrcpy is for __builtin___strcpy_chk.
// https://opensource.apple.com/source/Libc/Libc-498/include/secure/_string.h
func BuiltinStrcpy(dest, src []byte, size int) []byte {
	return noarch.Strcpy(dest, src)
}

// BuiltinObjectSize is for __builtin_object_size.
// https://github.com/elliotchance/c2go/issues/359
func BuiltinObjectSize(ptr []byte, theType int) int {
	return 5
}

// BuiltinStrncpy is for __builtin___strncpy_chk.
// https://opensource.apple.com/source/Libc/Libc-498/include/secure/_string.h
func BuiltinStrncpy(dest, src []byte, len, size int) []byte {
	return noarch.Strncpy(dest, src, len)
}

// BuiltinStrcat is for __builtin___strcat_chk
// https://opensource.apple.com/source/Libc/Libc-763.12/include/secure/_string.h.auto.html
func BuiltinStrcat(dest, src []byte, _ int) []byte {
	return noarch.Strcat(dest, src)
}
