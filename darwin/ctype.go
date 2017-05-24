package darwin

import (
	"unicode"
	"unicode/utf8"
)

// CtRuneT represents __darwin_ct_rune_t.
type CtRuneT int

// Apple defines a bunch of magic values for the type of character, see
// https://opensource.apple.com/source/Libc/Libc-320/include/ctype.h.auto.html
//
// These are provided as values for method that take _f parameter.
const (
	CTYPE_A   = 0x00000100 // Alpha
	CTYPE_C   = 0x00000200 // Control
	CTYPE_D   = 0x00000400 // Digit
	CTYPE_G   = 0x00000800 // Graph
	CTYPE_L   = 0x00001000 // Lower
	CTYPE_P   = 0x00002000 // Punct
	CTYPE_S   = 0x00004000 // Space
	CTYPE_U   = 0x00008000 // Upper
	CTYPE_X   = 0x00010000 // X digit
	CTYPE_B   = 0x00020000 // Blank
	CTYPE_R   = 0x00040000 // Print
	CTYPE_I   = 0x00080000 // Ideogram
	CTYPE_T   = 0x00100000 // Special
	CTYPE_Q   = 0x00200000 // Phonogram
	CTYPE_SW0 = 0x20000000 // 0 width character
	CTYPE_SW1 = 0x40000000 // 1 width character
	CTYPE_SW2 = 0x80000000 // 2 width character
	CTYPE_SW3 = 0xc0000000 // 3 width character
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
	if _f&CTYPE_A != 0 && unicode.IsLetter(rune(_c)) {
		return 1
	}

	if _f&CTYPE_C != 0 && unicode.IsControl(rune(_c)) {
		return 1
	}

	if _f&CTYPE_D != 0 && unicode.IsDigit(rune(_c)) {
		return 1
	}

	if _f&CTYPE_G != 0 && unicode.IsGraphic(rune(_c)) {
		return 1
	}

	if _f&CTYPE_L != 0 && unicode.IsLower(rune(_c)) {
		return 1
	}

	if _f&CTYPE_P != 0 && unicode.IsPunct(rune(_c)) {
		return 1
	}

	if _f&CTYPE_S != 0 && unicode.IsSpace(rune(_c)) {
		return 1
	}

	if _f&CTYPE_U != 0 && unicode.IsUpper(rune(_c)) {
		return 1
	}

	if _f&CTYPE_R != 0 && unicode.IsPrint(rune(_c)) {
		return 1
	}

	// TODO: Is this really the right way to do this?
	if _f&CTYPE_X != 0 && (unicode.IsDigit(rune(_c)) ||
		(_c >= 'a' && _c <= 'f') ||
		(_c >= 'A' && _c <= 'F')) {
		return 1
	}

	// These are not supported, yet.
	if _f&CTYPE_B != 0 {
		panic("CTYPE_B is not supported")
	}

	if _f&CTYPE_I != 0 {
		panic("CTYPE_I is not supported")
	}

	if _f&CTYPE_T != 0 {
		panic("CTYPE_T is not supported")
	}

	if _f&CTYPE_Q != 0 {
		panic("CTYPE_Q is not supported")
	}

	// Extra rules around the character width (in bytes).
	_, size := utf8.DecodeLastRuneInString(string(_c))

	if _f&CTYPE_SW0 != 0 && size == 0 {
		return 1
	}
	if _f&CTYPE_SW1 != 0 && size == 1 {
		return 1
	}
	if _f&CTYPE_SW2 != 0 && size == 2 {
		return 1
	}
	if _f&CTYPE_SW3 != 0 && size == 3 {
		return 1
	}

	return 0
}

// I have no idea what MaskRune is supposed to do. It is provided internally by
// darwin.
func MaskRune(_c CtRuneT, _f uint32) CtRuneT {
	return _c
}

// __darwin_ct_rune_t __isctype(__darwin_ct_rune_t, uint32)
func IsCType(_c CtRuneT, _f uint32) CtRuneT {
	return CtRuneT(IsType(_c, _f))
}

func ToLower(_c CtRuneT) CtRuneT {
	return CtRuneT(unicode.ToLower(rune(_c)))
}

func ToUpper(_c CtRuneT) CtRuneT {
	return CtRuneT(unicode.ToUpper(rune(_c)))
}
