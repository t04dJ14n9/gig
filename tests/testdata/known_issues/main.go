package known_issues

// ─────────────────────────────────────────────────────────────────────────────
// KNOWN ISSUES
//
// Resolved bugs are documented in testdata/resolved_issue/main.go as
// regression tests. This file tracks bugs awaiting fixes.
//
// Resolution mapping:
//   Bug 1 → Resolved Issue 28 (sort named-type conversion)
//   Bug 2 → fixed in gentool (time.Duration DirectCall)
//   Bug 3 → Resolved Issue 29 (fmt.Stringer)
//   Bug 4 → Resolved Issue 30 (fmt.Sprintf %T)
//   Bug 5 → Resolved Issue 31 (fmt.Sprintf %v _gig_id)
//   Bug 6 → Resolved Issue 32 (int64/uint64 narrowing)
//   Bug 7 → Resolved Issue 33 (bytes.Buffer.Cap)
// ─────────────────────────────────────────────────────────────────────────────

import (
	"bytes"
	"encoding/json"
)

// Bug 8: Method dispatch type collision — json.Encoder.Encode vs xml.Encoder.Encode
//
// When a program uses both json.Encoder and xml.Encoder, the compiled program
// contains methods with the same name "Encode" on different receiver types.
// The method resolver picks xml.Encoder.Encode for json.Encoder.Encode calls,
// causing: panic: interface {} is *json.Encoder, not *xml.Encoder
//
// This is the same reflect.StructOf type identity issue that was partially
// fixed for other cases. The fix likely requires making the method cache
// key include the full receiver type, not just the method name.

// JsonEncodeBug8 tests json.NewEncoder.Encode call.
func JsonEncodeBug8() int {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(map[string]int{"y": 20})
	return buf.Len()
}
