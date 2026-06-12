package divergence_hunt169

import (
	"fmt"
	"net/url"
	"strings"
)

// ============================================================================
// Round 169: URL parsing patterns (replaces HTTP which is not available)
// ============================================================================

// URLParsing tests URL parsing
func URLParsing() string {
	rawURL := "https://example.com:8080/path?key=value&foo=bar"
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Sprintf("error=%v", err)
	}
	return fmt.Sprintf("scheme=%s,host=%s,path=%s", u.Scheme, u.Host, u.Path)
}

// URLQuery tests query parameter handling
func URLQuery() string {
	rawURL := "https://example.com?a=1&b=2&c=3"
	u, _ := url.Parse(rawURL)
	query := u.Query()
	return fmt.Sprintf("a=%s,b=%s,c=%s", query.Get("a"), query.Get("b"), query.Get("c"))
}

// URLString tests URL reconstruction
func URLString() string {
	u := &url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/path",
	}
	u.RawQuery = "key=value"
	return u.String()
}

// URLEscaping tests URL escaping
func URLEscaping() string {
	s := "hello world & more"
	escaped := url.QueryEscape(s)
	unescaped, _ := url.QueryUnescape(escaped)
	return fmt.Sprintf("escaped=%s,unescaped=%s", escaped, unescaped)
}

// URLPathEscaping tests path escaping
func URLPathEscaping() string {
	path := "/path with spaces/file.txt"
	escaped := url.PathEscape(path)
	unescaped, _ := url.PathUnescape(escaped)
	return fmt.Sprintf("escaped=%s,unescaped=%s", escaped, unescaped)
}

// URLUserInfo tests user info in URL
func URLUserInfo() string {
	u := &url.URL{
		Scheme: "https",
		User:   url.UserPassword("user", "pass"),
		Host:   "example.com",
		Path:   "/path",
	}
	return fmt.Sprintf("user=%s", u.User.Username())
}

// URLIsAbs tests if URL is absolute
func URLIsAbs() string {
	abs, _ := url.Parse("https://example.com/path")
	rel, _ := url.Parse("/path")
	return fmt.Sprintf("abs=%v,rel=%v", abs.IsAbs(), rel.IsAbs())
}

// URLHostname tests hostname extraction
func URLHostname() string {
	u, _ := url.Parse("https://example.com:8080/path")
	return fmt.Sprintf("hostname=%s,port=%s", u.Hostname(), u.Port())
}

// URLRequestURI tests request URI
func URLRequestURI() string {
	u, _ := url.Parse("https://example.com/path?query=1")
	return fmt.Sprintf("requesturi=%s", u.RequestURI())
}

// URLResolveReference tests URL resolution
func URLResolveReference() string {
	base, _ := url.Parse("https://example.com/a/b/")
	ref, _ := url.Parse("c")
	resolved := base.ResolveReference(ref)
	return fmt.Sprintf("resolved=%s", resolved.String())
}

// URLValues tests url.Values
func URLValues() string {
	v := url.Values{}
	v.Set("key", "value")
	v.Add("multi", "a")
	v.Add("multi", "b")
	return fmt.Sprintf("key=%s,multi=%s", v.Get("key"), strings.Join(v["multi"], ","))
}

// URLEncode tests encoding url.Values
func URLEncode() string {
	v := url.Values{}
	v.Set("a", "1")
	v.Set("b", "2")
	return v.Encode()
}

// URLParseQuery tests parsing query string
func URLParseQuery() string {
	query := "a=1&b=2&a=3"
	values, _ := url.ParseQuery(query)
	return fmt.Sprintf("a=%s,b=%s", strings.Join(values["a"], ","), values.Get("b"))
}
