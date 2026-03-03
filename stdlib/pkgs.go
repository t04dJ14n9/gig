// Package stdlib declares dependencies for the gig interpreter.
//
// This package lists the standard library packages that are pre-registered for use
// in interpreted Go code. When you import gig/stdlib/packages, all these packages
// become available to your interpreted programs.
//
// # Included Packages
//
// The stdlib includes 39 standard library packages:
//
//   - bytes, strings, strconv: String and byte manipulation
//   - fmt: Formatted I/O
//   - math, math/rand/v2: Mathematical operations
//   - time: Time and duration handling
//   - encoding/json, encoding/xml, encoding/base64, encoding/csv, encoding/hex: Data encoding
//   - errors: Error handling
//   - sort, slices, maps, cmp: Collection utilities
//   - regexp: Regular expressions
//   - container/heap, container/list, container/ring: Container types
//   - context: Context handling
//   - crypto/hmac, crypto/sha256: Cryptographic functions
//   - html, html/template, text/template: HTML and template processing
//   - io: I/O primitives
//   - log: Logging
//   - net/http, net/url: HTTP client and URL handling
//   - os: Operating system interface
//   - path, path/filepath: Path manipulation
//   - sync, sync/atomic: Synchronization primitives
//   - unicode, unicode/utf8, unicode/utf16: Unicode handling
//
// # Usage
//
// Import the stdlib packages before building interpreted code:
//
//	import _ "git.woa.com/youngjin/gig/stdlib/packages"
//
//	func main() {
//	    prog, err := gig.Build(`
//	        package main
//	        import "fmt"
//	        func Hello() { fmt.Println("Hello, World!") }
//	    `)
//	    // ...
//	}
//
// # Adding Custom Packages
//
// To add third-party libraries or additional standard library packages:
//  1. Create a custom dependency package using `gig init -package mydep`
//  2. Edit mydep/pkgs.go to add your imports
//  3. Run `gig gen ./mydep`
//  4. Import both stdlib and your custom packages in your program
package stdlib

import (
	// ============================================
	// Go Standard Library (provided by gig)
	// ============================================
	_ "bytes"
	_ "cmp"
	_ "container/heap"
	_ "container/list"
	_ "container/ring"
	_ "context"
	_ "crypto/hmac"
	_ "crypto/sha256"
	_ "encoding/base64"
	_ "encoding/csv"
	_ "encoding/hex"
	_ "encoding/json"
	_ "encoding/xml"
	_ "errors"
	_ "fmt"
	_ "html"
	_ "html/template"
	_ "io"
	_ "log"
	_ "maps"
	_ "math"
	_ "math/rand/v2"
	_ "net/http"
	_ "net/url"
	_ "os"
	_ "path"
	_ "path/filepath"
	_ "regexp"
	_ "slices"
	_ "sort"
	_ "strconv"
	_ "strings"
	_ "sync"
	_ "sync/atomic"
	_ "text/template"
	_ "time"
	_ "unicode"
	_ "unicode/utf16"
	_ "unicode/utf8"
)
