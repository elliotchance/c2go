package darwin

import (
	"unicode"
	"unicode/utf8"
)

// CtRuneT represents __darwin_ct_rune_t.
type CtRuneT int32

// Apple defines a bunch of magic values for the type of character, see
// https://opensource.apple.com/source/Libc/Libc-320/include/ctype.h.auto.html
//
// These are provided as values for method that take _f parameter.
const (
	CtypeA   = 0x00000100 // Alpha
	CtypeC   = 0x00000200 // Control
	CtypeD   = 0x00000400 // Digit
	CtypeG   = 0x00000800 // Graph
	CtypeL   = 0x00001000 // Lower
	CtypeP   = 0x00002000 // Punct
	CtypeS   = 0x00004000 // Space
	CtypeU   = 0x00008000 // Upper
	CtypeX   = 0x00010000 // X digit
	CtypeB   = 0x00020000 // Blank
	CtypeR   = 0x00040000 // Print
	CtypeI   = 0x00080000 // Ideogram
	CtypeT   = 0x00100000 // Special
	CtypeQ   = 0x00200000 // Phonogram
	CtypeSW0 = 0x20000000 // 0 width character
	CtypeSW1 = 0x40000000 // 1 width character
	CtypeSW2 = 0x80000000 // 2 width character
	CtypeSW3 = 0xc0000000 // 3 width character
)

// IsType replaces __istype(). It should not be strictly necessary but the real
// __istype() refers to internal darwin state (_DefaultRuneLocale) that is
// difficult to translate. So for now we will replace it but this could be
// removed in the future.
//
// There may be multiple bit masks. And yes, I'm sure there is a much better way
// to handle this, so if you know one please consider putting in a PR :)
func IsType(_c CtRuneT, _f uint32) uint32 {
	// These are the easy ones.
	if _f&CtypeA != 0 && unicode.IsLetter(rune(_c)) && rune(_c) < 0x80 {
		return 1
	}

	if _f&CtypeC != 0 && unicode.IsControl(rune(_c)) && rune(_c) < 0x80 {
		return 1
	}

	if _f&CtypeD != 0 && unicode.IsDigit(rune(_c)) && rune(_c) < 0x80 {
		return 1
	}

	// The IsSpace check is required because Go treats spaces as graphic
	// characters, which C does not.
	if _f&CtypeG != 0 && unicode.IsGraphic(rune(_c)) && !unicode.IsSpace(rune(_c)) && rune(_c) < 0x80 {
		return 1
	}

	if _f&CtypeL != 0 && unicode.IsLower(rune(_c)) && rune(_c) < 0x80 {
		return 1
	}

	// Need to check for 0x24, 0x2b, 0x3c-0x3e, 0x5e, 0x60, 0x7c, 0x7e
	// because Go doesn't treat $+<=>^`|~ as punctuation.
	if _f&CtypeP != 0 && rune(_c) < 0x80 && (unicode.IsPunct(rune(_c)) || rune(_c) == 0x24 || rune(_c) == 0x2b ||
		(rune(_c) >= 0x3c && rune(_c) <= 0x3e) || rune(_c) == 0x5e || rune(_c) == 0x60 ||
		rune(_c) == 0x7c || rune(_c) == 0x7e) {
		return 1
	}

	if _f&CtypeS != 0 && unicode.IsSpace(rune(_c)) && rune(_c) < 0x80 {
		return 1
	}

	if _f&CtypeU != 0 && unicode.IsUpper(rune(_c)) && rune(_c) < 0x80 {
		return 1
	}

	if _f&CtypeR != 0 && unicode.IsPrint(rune(_c)) && rune(_c) < 0x80 {
		return 1
	}

	// TODO: Is this really the right way to do this?
	if _f&CtypeX != 0 && (unicode.IsDigit(rune(_c)) ||
		(_c >= 'a' && _c <= 'f') ||
		(_c >= 'A' && _c <= 'F')) && rune(_c) < 0x80 {
		return 1
	}

	// These are not supported, yet.
	if _f&CtypeB != 0 {
		panic("CtypeB is not supported")
	}

	if _f&CtypeI != 0 {
		panic("CtypeI is not supported")
	}

	if _f&CtypeT != 0 {
		panic("CtypeT is not supported")
	}

	if _f&CtypeQ != 0 {
		panic("CtypeQ is not supported")
	}

	// Extra rules around the character width (in bytes).
	_, size := utf8.DecodeLastRuneInString(string(_c))

	if _f&CtypeSW0 != 0 && size == 0 {
		return 1
	}
	if _f&CtypeSW1 != 0 && size == 1 {
		return 1
	}
	if _f&CtypeSW2 != 0 && size == 2 {
		return 1
	}
	if _f&CtypeSW3 != 0 && size == 3 {
		return 1
	}

	return 0
}

// MaskRune handles __maskrune(). I have no idea what MaskRune is supposed to
// do. It is provided internally by darwin.
func MaskRune(_c CtRuneT, _f uint32) CtRuneT {
	return _c
}

// IsCType handles __isctype.
func IsCType(_c CtRuneT, _f uint32) CtRuneT {
	return CtRuneT(IsType(_c, _f))
}

// ToLower handles __tolower().
func ToLower(_c CtRuneT) CtRuneT {
	return CtRuneT(unicode.ToLower(rune(_c)))
}

// ToUpper handles __toupper().
func ToUpper(_c CtRuneT) CtRuneT {
	return CtRuneT(unicode.ToUpper(rune(_c)))
}
