// Package tests v2 skip helper.
//
// When GIG_BACKEND=v2 is set, tests that exercise the legacy bytecode
// VM, the legacy WithStatefulGlobals mode, the VM pool, or DirectCall
// wrapper performance no longer apply: those subsystems are removed
// or unreachable from the new SSA pipeline. Each Category-B test file
// calls skipIfV2(t) at entry; the rest of the file is unmodified.
package tests

import (
	"os"
	"strings"
	"testing"
)

// isV2Backend reports whether the v2 SSA pipeline is active. v2 is now
// the default; only the explicit "legacy"/"v1" override falls back.
func isV2Backend() bool {
	switch strings.ToLower(os.Getenv("GIG_BACKEND")) {
	case "legacy", "v1", "vm":
		return false
	}
	return true
}

// skipIfV2 skips a test that targets the legacy backend's
// implementation details and is not applicable under v2.
func skipIfV2(t *testing.T, reason string) {
	t.Helper()
	if isV2Backend() {
		t.Skipf("skipped under GIG_BACKEND=v2: %s", reason)
	}
}

// toInt64 normalises a Run() result to an int64 for comparison.
// It accepts every integer width Go can produce.
func toInt64(v any) int64 {
	switch x := v.(type) {
	case int:
		return int64(x)
	case int8:
		return int64(x)
	case int16:
		return int64(x)
	case int32:
		return int64(x)
	case int64:
		return x
	case uint:
		return int64(x)
	case uint8:
		return int64(x)
	case uint16:
		return int64(x)
	case uint32:
		return int64(x)
	case uint64:
		return int64(x)
	}
	return 0
}

