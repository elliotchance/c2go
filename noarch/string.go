package noarch

import (
	"bytes"
	"reflect"
	"strings"
	"unsafe"
)

// Strlen returns the length of a string.
//
// The length of a C string is determined by the terminating null-character: A
// C string is as long as the number of characters between the beginning of the
// string and the terminating null character (without including the terminating
// null character itself).
func Strlen(a []byte) int32 {
	// TODO: The transpiler should have a syntax that means this proxy function
	// does not need to exist.

	return int32(len(CStringToString(a)))
}

// Strcpy copies the C string pointed by source into the array pointed by
// destination, including the terminating null character (and stopping at that
// point).
//
// To avoid overflows, the size of the array pointed by destination shall be
// long enough to contain the same C string as source (including the terminating
// null character), and should not overlap in memory with source.
func Strcpy(dest, src []byte) []byte {
	for i, c := range src {
		dest[i] = c

		// We only need to copy until the first NULL byte. Make sure we also
		// include that NULL byte on the end.
		if c == '\x00' {
			break
		}
	}

	return dest
}

// Strncpy copies the first num characters of source to destination. If the end
// of the source C string (which is signaled by a null-character) is found
// before num characters have been copied, destination is padded with zeros
// until a total of num characters have been written to it.
//
// No null-character is implicitly appended at the end of destination if source
// is longer than num. Thus, in this case, destination shall not be considered a
// null terminated C string (reading it as such would overflow).
//
// destination and source shall not overlap (see memmove for a safer alternative
// when overlapping).
func Strncpy(dest, src []byte, len int32) []byte {
	// Copy up to the len or first NULL bytes - whichever comes first.
	var i int32
	for ; i < len && src[i] != 0; i++ {
		dest[i] = src[i]
	}

	// The rest of the dest will be padded with zeros to the len.
	for ; i < len; i++ {
		dest[i] = 0
	}

	return dest
}

// Strcasestr - function is similar to Strstr(),
// but ignores the case of both strings.
func Strcasestr(str1, str2 []byte) []byte {
	a := strings.ToLower(CStringToString(str1))
	b := strings.ToLower(CStringToString(str2))
	index := strings.Index(a, b)
	if index == -1 {
		return nil
	}
	return str1[index : index+len(b)]
}

// Strcat - concatenate strings
// Appends a copy of the source string to the destination string.
// The terminating null character in destination is overwritten by the first
// character of source, and a null-character is included at the end
// of the new string formed by the concatenation of both in destination.
func Strcat(dest, src []byte) []byte {
	Strcpy(dest[Strlen(dest):], src)
	return dest
}

// Strcmp - compare two strings
// Compares the C string str1 to the C string str2.
func Strcmp(str1, str2 []byte) int32 {
	return int32(bytes.Compare([]byte(CStringToString(str1)), []byte(CStringToString(str2))))
}

// Strncmp - compare two strings
// Compares the C string str1 to the C string str2 upto the first NULL character
// or n-th character whichever comes first.
func Strncmp(str1, str2 []byte, n int32) int32 {
	a := []byte(CStringToString(str1))
	a = a[:int(min(int(n), len(a)))]
	b := []byte(CStringToString(str2))
	b = b[:int(min(int(n), len(b)))]
	return int32(bytes.Compare(a, b))
}

// Strstr - locate a substring in a string
// function locates the first occurrence of the null-terminated string needle
// in the null-terminated string haystack.
func Strstr(str1, str2 []byte) []byte {
	a := CStringToString(str1)
	b := CStringToString(str2)
	index := strings.Index(a, b)
	if index == -1 {
		return nil
	}
	return str1[index : index+len(b)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Strchr - Locate first occurrence of character in string
// See: http://www.cplusplus.com/reference/cstring/strchr/
func Strchr(str []byte, ch int32) []byte {
	i := 0
	for {
		if str[i] == '\x00' {
			break
		}
		if int32(str[i]) == ch {
			return str[i:]
		}
		i++
	}
	return nil
}

// Memset treats dst as a binary array and sets size bytes to the value val.
// Returns dst.
func Memset(dst interface{}, val int32, size int32) interface{} {
	vDst := reflect.ValueOf(dst).Type()
	switch vDst.Kind() {
	case reflect.Slice, reflect.Array:
		vDst = vDst.Elem()
	}
	baseSizeDst := int32(vDst.Size())
	data := *(*[]byte)(unsafe.Pointer(UnsafeSliceToSlice(dst, baseSizeDst, int32(1))))
	var i int32
	var vb = byte(val)
	for i = 0; i < size; i++ {
		data[i] = vb
	}
	return dst
}

// Memcpy treats dst and src as binary arrays and copies size bytes from src to dst.
// Returns dst.
// While in C it it is undefined behavior to call memcpy with overlapping regions,
// in Go we rely on the built-in copy function, which has no such limitation.
// To copy overlapping regions in C memmove should be used, so we map that function
// to Memcpy as well.
func Memcpy(dst interface{}, src interface{}, size int32) interface{} {
	vDst := reflect.ValueOf(dst).Type()
	switch vDst.Kind() {
	case reflect.Slice, reflect.Array:
		vDst = vDst.Elem()
	}
	baseSizeDst := int32(vDst.Size())
	vSrc := reflect.ValueOf(src).Type()
	switch vSrc.Kind() {
	case reflect.Slice, reflect.Array:
		vSrc = vSrc.Elem()
	}
	baseSizeSrc := int32(vSrc.Size())
	bDst := *(*[]byte)(unsafe.Pointer(UnsafeSliceToSlice(dst, baseSizeDst, int32(1))))
	bSrc := *(*[]byte)(unsafe.Pointer(UnsafeSliceToSlice(src, baseSizeSrc, int32(1))))
	copy(bDst[:size], bSrc[:size])
	return dst
}

// Memcmp compares two binary arrays upto n bytes.
// Different from strncmp, memcmp does not stop at the first NULL byte.
func Memcmp(src1, src2 interface{}, n int32) int32 {
	v1 := reflect.ValueOf(src1).Type()
	switch v1.Kind() {
	case reflect.Slice, reflect.Array:
		v1 = v1.Elem()
	}
	baseSize1 := int32(v1.Size())
	v2 := reflect.ValueOf(src2).Type()
	switch v2.Kind() {
	case reflect.Slice, reflect.Array:
		v2 = v2.Elem()
	}
	baseSize2 := int32(v2.Size())
	b1 := *(*[]byte)(unsafe.Pointer(UnsafeSliceToSlice(src1, baseSize1, int32(1))))
	b2 := *(*[]byte)(unsafe.Pointer(UnsafeSliceToSlice(src2, baseSize2, int32(1))))
	return int32(bytes.Compare(b1[:int(n)], b2[:int(n)]))
}
