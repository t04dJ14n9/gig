// Package mydep declares dependencies for gig interpreter.
// Standard library packages are included by default.
// Add your custom third-party library imports at the end.
package mydep

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

	// ============================================
	// Custom third-party libraries (add yours below)
	// ============================================
	_ "github.com/dromara/carbon/v2"
	_ "github.com/spf13/cast"
	_ "github.com/tidwall/gjson"
	_ "github.com/tidwall/sjson"
)
