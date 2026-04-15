package known_issues

// All previously known bugs have been resolved!
// See testdata/resolved_issue/main.go for regression tests.
//
// Resolved bug mapping:
//   Bug 1 → Resolved Issue 28 (sort named-type conversion)
//   Bug 2 → fixed in gentool (time.Duration DirectCall)
//   Bug 3 → Resolved Issue 29 (fmt.Stringer)
//   Bug 4 → Resolved Issue 30 (fmt.Sprintf %T)
//   Bug 5 → Resolved Issue 31 (fmt.Sprintf %v _gig_id)
//   Bug 6 → Resolved Issue 32 (int64/uint64 narrowing)
//   Bug 7 → Resolved Issue 33 (bytes.Buffer.Cap)
//   Bug 8 → Resolved Issue 34 (json.Encoder method dispatch collision)
//   StrangeSyntax Bug 1 → Resolved Issue 36 (nil slice to interface)
//   StrangeSyntax Bug 2 → Resolved Issue 37 (nil map access zero value)
//   StrangeSyntax Bug 3 → Resolved Issue 38 (nil map delete no-op)
//   StrangeSyntax Bug 4 → Resolved Issue 39 (blank expr interface return)
//   StrangeSyntax Bug 5 → Resolved Issue 40 (closed channel send recoverable)
//   StrangeSyntax Bug 6 → Resolved Issue 41 (nil func return type info)
//   PanicRecover bugs → All resolved (see panic_recover tests)
