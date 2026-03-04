package peephole

import "git.woa.com/youngjin/gig/bytecode"

// Pattern is the interface every peephole rule must implement.
// Match is called at position i in code; it returns the number of bytes consumed
// and the replacement bytes when the pattern fires, or ok=false to skip.
// A nil newBytes with ok=true means the matched bytes should be deleted entirely.
type Pattern interface {
	Match(code []byte, i int) (consumed int, newBytes []byte, ok bool)
}

// patterns is the ordered global registry of all peephole rules.
// Longer (more specific) patterns must be registered before shorter ones.
var patterns []Pattern

// Register appends one or more patterns to the global registry.
// It is intended to be called from init() functions in each pattern file.
func Register(p ...Pattern) {
	patterns = append(patterns, p...)
}

// Patterns returns the globally registered slice in registration order.
func Patterns() []Pattern {
	return patterns
}

// ---- shared builder helpers ----

// MatchOp reads the opcode at offset off, returning false if out of bounds.
func MatchOp(code []byte, off int, op bytecode.OpCode) bool {
	return off < len(code) && bytecode.OpCode(code[off]) == op
}

// Make3Op builds a 7-byte fused instruction: opcode + u16 + u16 + u16.
func Make3Op(op bytecode.OpCode, a, b, c uint16) []byte {
	out := make([]byte, 7)
	out[0] = byte(op)
	WriteU16(out, 1, a)
	WriteU16(out, 3, b)
	WriteU16(out, 5, c)
	return out
}

// Make2Op builds a 5-byte fused instruction: opcode + u16 + u16.
func Make2Op(op bytecode.OpCode, a, b uint16) []byte {
	out := make([]byte, 5)
	out[0] = byte(op)
	WriteU16(out, 1, a)
	WriteU16(out, 3, b)
	return out
}

// Make1Op builds a 3-byte fused instruction: opcode + u16.
func Make1Op(op bytecode.OpCode, a uint16) []byte {
	out := make([]byte, 3)
	out[0] = byte(op)
	WriteU16(out, 1, a)
	return out
}

// ReadU16 reads a big-endian uint16 from code at offset.
func ReadU16(code []byte, offset int) uint16 {
	return uint16(code[offset])<<8 | uint16(code[offset+1])
}

// WriteU16 writes a big-endian uint16 to code at offset.
func WriteU16(code []byte, offset int, val uint16) {
	code[offset] = byte(val >> 8)
	code[offset+1] = byte(val)
}
