// Package noarch contains low-level functions that apply to multiple platforms.
package noarch

const ISupper uint16 = ((1 << 0) << 8)
const ISlower uint16 = ((1 << 1) << 8)
const ISalpha uint16 = ((1 << 2) << 8)
const ISdigit uint16 = ((1 << 3) << 8)
const ISxdigit uint16 = ((1 << 4) << 8)
const ISspace uint16 = ((1 << 5) << 8)
const ISprint uint16 = ((1 << 6) << 8)
const ISgraph uint16 = ((1 << 7) << 8)
const ISblank uint16 = ((1 << 8) >> 8)
const IScntrl uint16 = ((1 << 9) >> 8)
const ISpunct uint16 = ((1 << 10) >> 8)
const ISalnum uint16 = ((1 << 11) >> 8)
