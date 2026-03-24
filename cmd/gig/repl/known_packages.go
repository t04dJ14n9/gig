// Package repl provides the core REPL session management.
package repl

// knownPackages maps package names to their import paths.
// This is used for auto-import detection when users reference packages
// without explicit import statements (e.g., fmt.Println).
var knownPackages = map[string]string{
	// Standard I/O & formatting
	"fmt": "fmt",
	"io":  "io",
	"log": "log",
	"os":  "os",

	// Strings & text
	"strings":      "strings",
	"strconv":      "strconv",
	"bytes":        "bytes",
	"regexp":       "regexp",
	"unicode":      "unicode",
	"utf8":         "unicode/utf8",
	"utf16":        "unicode/utf16",
	"template":     "text/template",
	"html":         "html",
	"htmltemplate": "html/template",

	// Encoding
	"json":   "encoding/json",
	"xml":    "encoding/xml",
	"base64": "encoding/base64",
	"hex":    "encoding/hex",
	"csv":    "encoding/csv",

	// Math & sorting
	"math": "math",
	"rand": "math/rand/v2",
	"sort": "sort",
	"cmp":  "cmp",

	// Collections
	"slices": "slices",
	"maps":   "maps",
	"heap":   "container/heap",
	"list":   "container/list",
	"ring":   "container/ring",

	// Net & URL
	"url":  "net/url",
	"http": "net/http",

	// Paths
	"path":     "path",
	"filepath": "path/filepath",

	// Time & context
	"time":    "time",
	"context": "context",

	// Concurrency
	"sync":   "sync",
	"atomic": "sync/atomic",
	"errors": "errors",

	// Crypto — standard library
	"crypto":     "crypto",
	"hmac":       "crypto/hmac",
	"sha256":     "crypto/sha256",
	"sha512":     "crypto/sha512",
	"md5":        "crypto/md5",
	"sha1":       "crypto/sha1",
	"aes":        "crypto/aes",
	"cipher":     "crypto/cipher",
	"des":        "crypto/des",
	"dsa":        "crypto/dsa",
	"ecdsa":      "crypto/ecdsa",
	"ed25519":    "crypto/ed25519",
	"elliptic":   "crypto/elliptic",
	"randcrypto": "crypto/rand",
	"rsa":        "crypto/rsa",
	"subtle":     "crypto/subtle",
	"tls":        "crypto/tls",
	"x509":       "crypto/x509",
	"pkix":       "crypto/x509/pkix",

	// Crypto — extended (golang.org/x/crypto)
	"bcrypt":  "golang.org/x/crypto/bcrypt",
	"blake2b": "golang.org/x/crypto/blake2b",
	"blake2s": "golang.org/x/crypto/blake2s",
	"scrypt":  "golang.org/x/crypto/scrypt",
	"argon2":  "golang.org/x/crypto/argon2",
}
