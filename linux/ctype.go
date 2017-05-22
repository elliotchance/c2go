package linux

import (
	"unicode"
)

var characterTable []uint16

func generateCharacterTable() {
	for i := 0; i < 255; i++ {
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

		if unicode.IsGraphic(rune(i)) {
			c |= ((1 << (7)) << 8)
		}

		// FIXME: Blank is not implemented.
		// if unicode.IsBlank(rune(i)) {
		// 	c |= ((1 << (8)) >> 8)
		// }

		if unicode.IsControl(rune(i)) {
			c |= ((1 << (9)) >> 8)
		}

		if unicode.IsPunct(rune(i)) {
			c |= ((1 << (10)) >> 8)
		}

		if unicode.IsLetter(rune(i)) || unicode.IsDigit(rune(i)) {
			c |= ((1 << (11)) >> 8)
		}

		// Yes, I know this is a hideously slow way to do it but I just want to
		// test if this works right now.
		characterTable = append(characterTable, c)
	}
}

func CtypeLoc() [][]uint16 {
	if len(characterTable) == 0 {
		generateCharacterTable()
	}

	return [][]uint16{characterTable}
}

func ToLower(_c int) int {
	return int(unicode.ToLower(rune(_c)))
}

func ToUpper(_c int) int {
	return int(unicode.ToUpper(rune(_c)))
}
