// Package stdlib declares dependencies for the gig interpreter.
//
// This package lists the standard library packages that are pre-registered for use
// in interpreted Go code. When you import gig/stdlib/packages, all these packages
// become available to your interpreted programs.
//
// # Sandbox Safety
//
// Packages are categorized as SAFE (no host/external interaction) or UNSAFE
// (filesystem, network, process, or system access). Only SAFE packages are enabled
// by default. UNSAFE packages are commented out and can be selectively enabled
// by the embedding application when needed.
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
	// SAFE: Pure computation, no host interaction
	// ============================================

	// --- String & byte manipulation ---
	_ "bufio"
	_ "bytes"
	_ "strings"
	_ "strconv"
	_ "unicode"
	_ "unicode/utf8"
	_ "unicode/utf16"

	// --- Formatted I/O ---
	_ "fmt"

	// --- Math ---
	_ "math"
	_ "math/big"
	_ "math/bits"
	_ "math/cmplx"
	_ "math/rand"
	_ "math/rand/v2"

	// --- Collections & sorting ---
	_ "cmp"
	_ "slices"
	_ "maps"
	_ "sort"
	_ "container/heap"
	_ "container/list"
	_ "container/ring"

	// --- Encoding & serialization ---
	_ "encoding/ascii85"
	_ "encoding/asn1"
	_ "encoding/base32"
	_ "encoding/base64"
	_ "encoding/binary"
	_ "encoding/csv"
	_ "encoding/gob"
	_ "encoding/hex"
	_ "encoding/json"
	_ "encoding/pem"
	_ "encoding/xml"

	// --- Error handling ---
	_ "errors"

	// --- Regular expressions ---
	_ "regexp"
	_ "regexp/syntax"

	// --- HTML/text escaping ---
	_ "html"

	// --- Time ---
	_ "time"

	// --- Hashing (non-crypto, pure computation) ---
	_ "hash/adler32"
	_ "hash/crc32"
	_ "hash/crc64"
	_ "hash/fnv"
	_ "hash/maphash"

	// --- Cryptographic hashing (pure computation, no I/O) ---
	_ "crypto/hmac"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	_ "crypto/subtle"

	// --- Symmetric ciphers (pure computation, no I/O) ---
	_ "crypto/aes"
	_ "crypto/cipher"
	_ "crypto/des"
	_ "crypto/elliptic"
	_ "crypto/rc4"

	// --- Path manipulation (pure string ops, no filesystem) ---
	_ "path"

	// --- Concurrency (needed because gig supports `go` keyword) ---
	_ "sync"
	_ "context"

	// --- I/O interfaces (needed for encoding/json.NewDecoder etc.) ---
	_ "io"

	// --- Text processing ---
	_ "text/scanner"
	_ "text/tabwriter"

	// --- URL / IP parsing (pure string ops, no network) ---
	_ "net/url"
	_ "net/netip"

	// --- Archive (pure computation on byte streams) ---
	_ "archive/tar"
	_ "archive/zip"

	// --- Compression ---
	_ "compress/bzip2"
	_ "compress/flate"
	_ "compress/gzip"
	_ "compress/lzw"
	_ "compress/zlib"

	// --- MIME ---
	_ "mime"
	_ "mime/multipart"
	_ "mime/quotedprintable"

	// --- Image (color only, pure math) ---
	_ "image/color"

	// --- Misc pure computation ---
	_ "index/suffixarray"

	// ============================================
	// UNSAFE: Host/external interaction — disabled
	// ============================================

	// --- Filesystem ---
	// _ "os"             // file I/O, env vars, process exit
	// _ "os/exec"        // execute external commands
	// _ "os/signal"      // signal handling
	// _ "os/user"        // user account lookups
	// _ "path/filepath"  // OS-specific path ops (needs os)
	// _ "io/fs"          // filesystem abstraction

	// --- Network ---
	// _ "net"            // TCP/UDP sockets, DNS resolution
	// _ "net/http"       // HTTP client/server
	// _ "net/mail"       // mail parsing (safe, but niche)
	// _ "net/smtp"       // SMTP client
	// _ "net/rpc"        // RPC client/server
	// _ "net/textproto"  // text-based protocol support

	// --- Logging (writes to os.Stderr) ---
	// _ "log"
	// _ "log/slog"
	// _ "log/syslog"

	// --- Templates (can call arbitrary methods on objects) ---
	// _ "html/template"
	// _ "text/template"
	// _ "text/template/parse"

	// --- Concurrency (low-level, sync.Mutex is sufficient) ---
	// _ "sync/atomic"

	// --- System/runtime ---
	// _ "syscall"
	// _ "plugin"
	// _ "unsafe"
	// _ "runtime"
	// _ "runtime/cgo"
	// _ "runtime/debug"
	// _ "runtime/metrics"
	// _ "runtime/pprof"
	// _ "runtime/race"
	// _ "runtime/trace"
	// _ "embed"          // requires compiler support

	// --- Crypto (asymmetric — depends on crypto/rand for key gen) ---
	// _ "crypto/rand"    // reads from /dev/urandom
	// _ "crypto/ecdsa"   // depends on crypto/rand
	// _ "crypto/ed25519" // depends on crypto/rand
	// _ "crypto/rsa"     // depends on crypto/rand
	// _ "crypto/ecdh"    // depends on crypto/rand
	// _ "crypto/dsa"     // deprecated, depends on crypto/rand
	// _ "crypto/tls"     // network + filesystem
	// _ "crypto/x509"    // filesystem (cert store)

	// --- Database (needs external connection) ---
	// _ "database/sql"
	// _ "database/sql/driver"

	// --- Debug/binary inspection ---
	// _ "reflect"
	// _ "debug/buildinfo"
	// _ "debug/dwarf"
	// _ "debug/elf"
	// _ "debug/gosym"
	// _ "debug/macho"
	// _ "debug/pe"
	// _ "debug/plan9obj"

	// --- Testing ---
	// _ "testing"
	// _ "testing/fstest"
	// _ "testing/iotest"
	// _ "testing/quick"

	// --- Image (large dependency, rarely needed for rules) ---
	// _ "image"
	// _ "image/color/palette"
	// _ "image/draw"
	// _ "image/gif"
	// _ "image/jpeg"
	// _ "image/png"

	// --- Go toolchain (safe but very niche) ---
	// _ "go/ast"
	// _ "go/build"
	// _ "go/build/constraint"
	// _ "go/constant"
	// _ "go/doc"
	// _ "go/doc/comment"
	// _ "go/format"
	// _ "go/importer"
	// _ "go/parser"
	// _ "go/printer"
	// _ "go/scanner"
	// _ "go/token"
	// _ "go/types"
	// _ "go/version"

	// --- CLI ---
	// _ "flag"           // reads os.Args
	// _ "expvar"         // registers HTTP handlers

	// --- IO utilities ---
	// _ "io/ioutil"      // deprecated
)
