// Package darwin contains low-level functions for the Darwin (macOS) operating
// system.
package darwin

// FIXME: These are wrong.
type C__mbstate_t int64

// I'm not sure which header file this actually comes from?
func BSwap32(a uint32) uint32 {
	panic("BSwap32 is not supported")
}

// I'm not sure which header file this actually comes from?
func BSwap64(a uint64) uint64 {
	panic("BSwap64 is not supported")
}
