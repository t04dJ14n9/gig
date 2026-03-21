package thirdparty

import (
	"strings"
	"text/scanner"
	"text/tabwriter"
)

// ============================================================================
// text/scanner — lexical scanner
// ============================================================================

// TextScannerBasic tests basic token scanning.
func TextScannerBasic() int {
	var s scanner.Scanner
	s.Init(strings.NewReader("hello world 123"))
	var tok rune
	var count int
	for tok != rune(scanner.EOF) {
		tok = s.Scan()
		if tok != rune(scanner.EOF) {
			count++
		}
	}
	if count >= 3 {
		return 1
	}
	return 0
}

// TextScannerInts tests scanning integers.
func TextScannerInts() int {
	var s scanner.Scanner
	s.Init(strings.NewReader("42 100 -5 3.14"))
	s.Mode = uint(scanner.ScanInts)
	var count int
	for tok := s.Scan(); tok != rune(scanner.EOF); tok = s.Scan() {
		count++
	}
	if count >= 3 {
		return 1
	}
	return 0
}

// TextScannerStrings tests scanning string literals.
func TextScannerStrings() int {
	var s scanner.Scanner
	s.Init(strings.NewReader(`"hello" "world"`))
	s.Mode = uint(scanner.ScanStrings)
	var count int
	for tok := s.Scan(); tok != rune(scanner.EOF); tok = s.Scan() {
		count++
	}
	if count >= 2 {
		return 1
	}
	return 0
}

// TextScannerPosition tests position tracking.
func TextScannerPosition() int {
	var s scanner.Scanner
	s.Init(strings.NewReader("hello world"))
	s.Scan() // "hello"
	pos := s.Pos()
	if pos.Offset == 5 {
		return 1
	}
	return 0
}

// ============================================================================
// text/tabwriter — tab-aligned formatting
// ============================================================================

// TextTabwriterInit tests tabwriter initialization.
func TextTabwriterInit() int {
	w := tabwriter.NewWriter(&strings.Builder{}, 0, 0, 1, ' ', 0)
	if w != nil {
		return 1
	}
	return 0
}
