package optimize

import (
	"github.com/t04dJ14n9/gig/compiler/peephole"
	"github.com/t04dJ14n9/gig/model/bytecode"
)

// Peephole performs peephole optimization on compiled bytecode.
func Peephole(code []byte) []byte {
	var rewrites []rewrite

	i := 0
	for i < len(code) {
		op := bytecode.OpCode(code[i])
		instrEnd := i + 1 + opcodeWidth(op)
		if instrEnd > len(code) {
			break
		}

		matched := false
		for _, p := range peephole.Patterns() {
			if consumed, newBytes, ok := p.Match(code, i); ok {
				rewrites = append(rewrites, rewrite{i, i + consumed, newBytes})
				i += consumed
				matched = true
				break
			}
		}
		if !matched {
			i = instrEnd
		}
	}

	if len(rewrites) == 0 {
		return code
	}
	return applyRewrites(code, rewrites)
}
