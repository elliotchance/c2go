package linux

import (
	"unicode"
	"unsafe"
)

var characterTable []uint16

func generateCharacterTable() {
	for i := 0; i < 0x80; i++ {
		var c uint16

		// Each of the bitwise expressions below were copied from the enum
		// values, like _ISupper, etc.

		if unicode.IsUpper(rune(i)) {
			c |= ((1 << (0)) << 8)
		}

		if unicode.IsLower(rune(i)) {
			c |= ((1 << (1)) << 8)
		}

		if unicode.IsLetter(rune(i)) {
			c |= ((1 << (2)) << 8)
		}

		if unicode.IsDigit(rune(i)) {
			c |= ((1 << (3)) << 8)
		}

		if unicode.IsDigit(rune(i)) ||
			(i >= 'a' && i <= 'f') ||
			(i >= 'A' && i <= 'F') {
			// IsXDigit. This is the same implementation as the Mac version.
			// There may be a better way to do this.
			c |= ((1 << (4)) << 8)
		}

		if unicode.IsSpace(rune(i)) {
			c |= ((1 << (5)) << 8)
		}

		if unicode.IsPrint(rune(i)) {
			c |= ((1 << (6)) << 8)
		}

		// The IsSpace check is required because Go treats spaces as graphic
		// characters, which C does not.
		if unicode.IsGraphic(rune(i)) && !unicode.IsSpace(rune(i)) {
			c |= ((1 << (7)) << 8)
		}

		// FIXME: Blank is not implemented.
		// if unicode.IsBlank(rune(i)) {
		// 	c |= ((1 << (8)) >> 8)
		// }

		if unicode.IsControl(rune(i)) {
			c |= ((1 << (9)) >> 8)
		}

		// Need to check for 0x24, 0x2b, 0x3c-0x3e, 0x5e, 0x60, 0x7c, 0x7e
		// because Go doesn't treat $+<=>^`|~ as punctuation.
		if unicode.IsPunct(rune(i)) || i == 0x24 || i == 0x2b || (i >= 0x3c && i <= 0x3e) || i == 0x5e || i == 0x60 ||
			i == 0x7c || i == 0x7e {
			c |= ((1 << (10)) >> 8)
		}

		if unicode.IsLetter(rune(i)) || unicode.IsDigit(rune(i)) {
			c |= ((1 << (11)) >> 8)
		}

		// Yes, I know this is a hideously slow way to do it but I just want to
		// test if this works right now.
		characterTable = append(characterTable, c)
	}
	for i := 0x80; i < 256; i++ {
		// false for all characters > 0x7f
		characterTable = append(characterTable, 0)
	}
}

// CtypeLoc handles __ctype_b_loc(). It returns a character table.
func CtypeLoc() **uint16 {
	if len(characterTable) == 0 {
		generateCharacterTable()
	}

	return (**uint16)(unsafe.Pointer(&characterTable))
}

// ToLower handles tolower().
func ToLower(_c int32) int32 {
	return int32(unicode.ToLower(rune(_c)))
}

// ToUpper handles toupper().
func ToUpper(_c int32) int32 {
	return int32(unicode.ToUpper(rune(_c)))
}
