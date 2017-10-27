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

		// The IsSpace check is required becuase Go treats spaces as graphic
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

// CtypeLoc handles __ctype_b_loc(). It returns a character table.
func CtypeLoc() [][]uint16 {
	if len(characterTable) == 0 {
		generateCharacterTable()
	}

	return [][]uint16{characterTable}
}

const (
	cFalse int = 0
	cTrue  int = 1
)

func IsAlpha(_c int) int {
	if _c < 'A' || _c > 'z' {
		return cFalse
	} else if _c > 'Z' && _c < 'a' {
		return cFalse
	}
	return cTrue
}

func IsAlnum(_c int) int {
	if IsDigit(_c) == cTrue {
		return cTrue
	}
	if IsAlpha(_c) == cTrue {
		return cTrue
	}
	return cFalse
}

func IsCntrl(_c int) int {
	if unicode.IsControl(rune(_c)) {
		return cTrue
	}
	return cFalse
}

func IsDigit(_c int) int {
	if _c >= '0' && _c <= '9' {
		return cTrue
	}
	return cFalse
}

func IsGraph(_c int) int {
	if _c == ' ' {
		return cFalse // TODO : Check - some different between C and Go
	}
	if unicode.IsGraphic(rune(_c)) {
		return cTrue
	}
	return cFalse
}

func IsLower(_c int) int {
	if unicode.IsLower(rune(_c)) {
		return cTrue
	}
	return cFalse
}

func IsPrint(_c int) int {
	if unicode.IsPrint(rune(_c)) {
		return cTrue
	}
	return cFalse
}

func IsPunct(_c int) int {
	if unicode.IsPunct(rune(_c)) {
		return cTrue
	}
	return cFalse
}

func IsSpace(_c int) int {
	if unicode.IsSpace(rune(_c)) {
		return cTrue
	}
	return cFalse
}

func IsUpper(_c int) int {
	if unicode.IsUpper(rune(_c)) {
		return cTrue
	}
	return cFalse
}

func IsXDigit(_c int) int {
	if _c >= '0' && _c <= '9' {
		return cTrue
	}
	if _c >= 'A' && _c <= 'F' {
		return cTrue
	}
	if _c >= 'a' && _c <= 'f' {
		return cTrue
	}
	return cFalse
}

// ToUpper handles toupper().
func ToUpper(_c int) int {
	return int(unicode.ToUpper(rune(_c)))
}

// ToLower handles tolower().
func ToLower(_c int) int {
	return int(unicode.ToLower(rune(_c)))
}

func IsAscii(_c int) int {
	// TODO
	return cFalse
}

func ToAscii(_c int) int {
	// TODO
	return cFalse
}
