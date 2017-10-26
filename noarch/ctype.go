// Package noarch contains low-level functions that apply to multiple platforms.
package noarch

// Constants of `ctype.h`
const (
	ISupper  uint16 = ((1 << 0) << 8)
	ISlower  uint16 = ((1 << 1) << 8)
	ISalpha  uint16 = ((1 << 2) << 8)
	ISdigit  uint16 = ((1 << 3) << 8)
	ISxdigit uint16 = ((1 << 4) << 8)
	ISspace  uint16 = ((1 << 5) << 8)
	ISprint  uint16 = ((1 << 6) << 8)
	ISgraph  uint16 = ((1 << 7) << 8)
	ISblank  uint16 = ((1 << 8) >> 8)
	IScntrl  uint16 = ((1 << 9) >> 8)
	ISpunct  uint16 = ((1 << 10) >> 8)
	ISalnum  uint16 = ((1 << 11) >> 8)
)
