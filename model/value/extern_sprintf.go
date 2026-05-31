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
	// Fast path: no %T in format string — use standard fmt.Sprintf
	if !strings.Contains(format, "%T") {
		return fmt.Sprintf(format, args...)
	}
	// Slow path: replace %T for gigStructWrapper args with their type name
	var result strings.Builder
	argIdx := 0
	i := 0
	for i < len(format) {
		if format[i] == '%' {
			if i+1 < len(format) && format[i+1] == '%' {
				result.WriteString("%%")
				i += 2
				continue
			}
			j := i + 1
			// Skip flags
			for j < len(format) && (format[j] == '-' || format[j] == '+' || format[j] == '#' || format[j] == ' ' || format[j] == '0') {
				j++
			}
			// Skip width
			for j < len(format) && format[j] >= '0' && format[j] <= '9' {
				j++
			}
			// Skip precision
			if j < len(format) && format[j] == '.' {
				j++
				for j < len(format) && format[j] >= '0' && format[j] <= '9' {
					j++
				}
			}
			if j < len(format) {
				verb := format[j]
				if verb == 'T' && argIdx < len(args) {
					if w, ok := args[argIdx].(*gigStructWrapper); ok {
						result.WriteString(w.typeName)
						argIdx++
						i = j + 1
						continue
					}
				}
				if argIdx < len(args) {
					_, _ = fmt.Fprintf(&result, format[i:j+1], args[argIdx])
					argIdx++
				} else {
					result.WriteString(format[i : j+1])
				}
				i = j + 1
				continue
			}
		}
		result.WriteByte(format[i])
		i++
	}
	return result.String()
}
