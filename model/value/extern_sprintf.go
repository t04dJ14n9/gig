package value

import (
	"fmt"
	"strings"
)

// SprintfExtern is a general-purpose fmt.Sprintf replacement that correctly
// handles %T for gigStructWrapper values. Go's fmt.Sprintf("%T") bypasses
// fmt.Formatter entirely and uses reflect.TypeOf().String(), so we must
// intercept %T ourselves.
func SprintfExtern(format string, args ...any) string {
	if !strings.Contains(format, "%T") {
		return fmt.Sprintf(format, args...)
	}
	return formatSprintfExtern(format, args)
}

func formatSprintfExtern(format string, args []any) string {
	var result strings.Builder
	argIdx := 0
	pos := 0
	for pos < len(format) {
		pos, argIdx = writeNextSprintfExtern(&result, format, args, pos, argIdx)
	}
	return result.String()
}

func writeNextSprintfExtern(result *strings.Builder, format string, args []any, pos int, argIdx int) (int, int) {
	if format[pos] != '%' {
		return writeSprintfLiteralByte(result, format, pos), argIdx
	}
	if nextPos, ok := writeSprintfEscapedPercent(result, format, pos); ok {
		return nextPos, argIdx
	}
	if directive, ok := scanSprintfDirective(format, pos); ok {
		return writeSprintfDirective(result, format, args, directive, argIdx)
	}
	return writeSprintfLiteralByte(result, format, pos), argIdx
}

func writeSprintfLiteralByte(result *strings.Builder, format string, pos int) int {
	result.WriteByte(format[pos])
	return pos + 1
}

func writeSprintfEscapedPercent(result *strings.Builder, format string, pos int) (int, bool) {
	if pos+1 >= len(format) || format[pos+1] != '%' {
		return pos, false
	}
	// Preserve the legacy slow-path behavior. This path exists only for
	// formats containing %T, where Gig manually scans directives.
	result.WriteString("%%")
	return pos + 2, true
}

func writeSprintfDirective(
	result *strings.Builder,
	format string,
	args []any,
	d sprintfDirective,
	argIdx int,
) (int, int) {
	if nextArg, ok := writeGigTypeDirective(result, args, d, argIdx); ok {
		return d.end, nextArg
	}
	if argIdx < len(args) {
		_, _ = fmt.Fprintf(result, d.text(format), args[argIdx])
		argIdx++
	} else {
		result.WriteString(d.text(format))
	}
	return d.end, argIdx
}

func writeGigTypeDirective(
	result *strings.Builder,
	args []any,
	d sprintfDirective,
	argIdx int,
) (int, bool) {
	if d.verb != 'T' || argIdx >= len(args) {
		return argIdx, false
	}
	w, ok := args[argIdx].(*gigStructWrapper)
	if !ok {
		return argIdx, false
	}
	result.WriteString(w.typeName)
	return argIdx + 1, true
}

type sprintfDirective struct {
	start int
	end   int
	verb  byte
}

func (d sprintfDirective) text(format string) string {
	return format[d.start:d.end]
}

func scanSprintfDirective(format string, start int) (sprintfDirective, bool) {
	j := skipSprintfFlags(format, start+1)
	j = skipSprintfDigits(format, j)
	j = skipSprintfPrecision(format, j)
	if j >= len(format) {
		return sprintfDirective{}, false
	}
	return sprintfDirective{start: start, end: j + 1, verb: format[j]}, true
}

func skipSprintfFlags(format string, pos int) int {
	for pos < len(format) && isSprintfFlag(format[pos]) {
		pos++
	}
	return pos
}

func isSprintfFlag(ch byte) bool {
	return ch == '-' || ch == '+' || ch == '#' || ch == ' ' || ch == '0'
}

func skipSprintfDigits(format string, pos int) int {
	for pos < len(format) && isSprintfDigit(format[pos]) {
		pos++
	}
	return pos
}

func skipSprintfPrecision(format string, pos int) int {
	if pos >= len(format) || format[pos] != '.' {
		return pos
	}
	return skipSprintfDigits(format, pos+1)
}

func isSprintfDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
