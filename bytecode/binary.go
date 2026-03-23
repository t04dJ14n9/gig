package bytecode

// ReadU16 reads a big-endian uint16 from code at offset.
func ReadU16(code []byte, offset int) uint16 {
	return uint16(code[offset])<<8 | uint16(code[offset+1])
}

// WriteU16 writes a big-endian uint16 to code at offset.
func WriteU16(code []byte, offset int, val uint16) {
	code[offset] = byte(val >> 8)
	code[offset+1] = byte(val)
}
