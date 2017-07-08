// Package darwin contains low-level functions for the Darwin (macOS) operating
// system.
package darwin

// BSwap32 handles __builtin_bswap32(). It is not supported and if used with
// panic. The original documentation says:
//
// Returns x with the order of the bytes reversed; for example, 0xaabb becomes
// 0xbbaa. Byte here always means exactly 8 bits.
func BSwap32(a uint32) uint32 {
	panic("BSwap32 is not supported")
}

// BSwap64 handles __builtin_bswap64(). It is not supported, see BSwap32().
func BSwap64(a uint64) uint64 {
	panic("BSwap64 is not supported")
}
