// Package pkgs provides a boilerplate pkgs.go template.
//
// Copy the BoilerplatePkgsGo string to create your own pkgs.go file,
// then add custom third-party imports at the end.
package pkgs

// BoilerplatePkgsGo is the default pkgs.go template that users can copy.
// It includes all Go standard library packages supported by gig.
// Users should add their third-party library imports after the standard library section.
const BoilerplatePkgsGo = `// Package pkgs contains gig interpreter dependencies.
// Standard library packages are included by default.
// Add your custom third-party library imports at the end.
package pkgs

import (
	// ============================================
	// Go Standard Library (default, provided by gig)
	// ============================================
	_ "bytes"
	_ "cmp"
	_ "container/heap"
	_ "container/list"
	_ "container/ring"
	_ "context"
	_ "crypto/hmac"
	_ "crypto/md5"
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
	_ "math/rand"
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
	_ "unicode/utf8"
	_ "unicode/utf16"

	// ============================================
	// Custom third-party libraries (add yours below)
	// ============================================
	// _ "github.com/spf13/cast"
	// _ "github.com/tidwall/gjson"
	// _ "github.com/tidwall/sjson"
)
`
