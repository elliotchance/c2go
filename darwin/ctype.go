package darwin

import (
	"unicode"
	"github.com/elliotchance/c2go/noarch"
)

// __darwin_ct_rune_t
type Darwin_ct_rune_t int

// IsAscii is properly provided, but it needs to be referenced by IsType.
func IsAscii(c Darwin_ct_rune_t) bool {
	return c <= unicode.MaxASCII
}

// IsType replaces __istype(). It should not be strictly necessary but the real
// __istype() refers to internal darwin state (_DefaultRuneLocale) that is
// difficult to translate. So for now we will replace it but this could be
// removed in the future.
func IsType(_c Darwin_ct_rune_t, _f uint32) uint32 {
	if IsAscii(_c) {
		return uint32(noarch.BoolToInt(unicode.IsLetter(rune(_c)) || unicode.IsDigit(rune(_c))))
	}

	return 0 // uint32(MaskRune(_c, _f))
}

// I have no idea what MaskRune is supposed to do. It is provided internally by
// darwin.
func MaskRune(_c Darwin_ct_rune_t, _f uint32) Darwin_ct_rune_t {
	return _c
}

// __darwin_ct_rune_t __isctype(__darwin_ct_rune_t, uint32)
func IsCType(_c Darwin_ct_rune_t, _f uint32) Darwin_ct_rune_t {
	if (_c < 0 || _c >= (1 <<8 )) {
		return 0
	}

	return _c & Darwin_ct_rune_t(_f);
}

func ToLower(_c Darwin_ct_rune_t) Darwin_ct_rune_t {
	return Darwin_ct_rune_t(unicode.ToLower(rune(_c)))
}

func ToUpper(_c Darwin_ct_rune_t) Darwin_ct_rune_t {
	return Darwin_ct_rune_t(unicode.ToUpper(rune(_c)))
}
