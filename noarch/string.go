package noarch

// Strlen returns the length of a string.
//
// The length of a C string is determined by the terminating null-character: A
// C string is as long as the number of characters between the beginning of the
// string and the terminating null character (without including the terminating
// null character itself).
func Strlen(a []byte) int {
	// TODO: The transpiler should have a syntax that means this proxy function
	// does not need to exist.

	return len(CStringToString(a))
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
func Strncpy(dest, src []byte, len int) []byte {
	// Copy up to the len or first NULL bytes - whichever comes first.
	i := 0
	for ; i < len && src[i] != 0; i++ {
		dest[i] = src[i]
	}

	// The rest of the dest will be padded with zeros to the len.
	for ; i < len; i++ {
		dest[i] = 0
	}

	return dest
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
