package darwin

import (
	"github.com/elliotchance/c2go/noarch"
	"unsafe"
)

// BuiltinStrcpy is for __builtin___strcpy_chk.
// https://opensource.apple.com/source/Libc/Libc-498/include/secure/_string.h
func BuiltinStrcpy(dest, src *byte, size int32) *byte {
	return noarch.Strcpy(dest, src)
}

// BuiltinObjectSize is for __builtin_object_size.
// https://github.com/elliotchance/c2go/issues/359
func BuiltinObjectSize(ptr *byte, theType int32) int32 {
	return 5
}

// BuiltinStrncpy is for __builtin___strncpy_chk.
// https://opensource.apple.com/source/Libc/Libc-498/include/secure/_string.h
func BuiltinStrncpy(dest, src *byte, len, size int32) *byte {
	return noarch.Strncpy(dest, src, len)
}

// BuiltinStrcat is for __builtin___strcat_chk
// https://opensource.apple.com/source/Libc/Libc-763.12/include/secure/_string.h.auto.html
func BuiltinStrcat(dest, src *byte, _ int32) *byte {
	return noarch.Strcat(dest, src)
}

// Memset is for __builtin___memset_chk
// https://opensource.apple.com/source/Libc/Libc-498/include/secure/_string.h
func Memset(dst unsafe.Pointer, val int32, size int32, _ int32) unsafe.Pointer {
	return noarch.Memset(dst, val, size)
}

// Memcpy  is for __builtin___memcpy_chk and __builtin___memmove_chk
//// https://opensource.apple.com/source/Libc/Libc-498/include/secure/_string.h
func Memcpy(dst, src unsafe.Pointer, size int32, _ int32) unsafe.Pointer {
	return noarch.Memcpy(dst, src, size)
}
