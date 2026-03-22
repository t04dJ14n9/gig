package thirdparty

import (
	"net/netip"
	"net/url"
)

// ============================================================================
// net/url — URL parsing and construction
// ============================================================================

// NetURLParse tests URL parsing.
func NetURLParse() int {
	u, err := url.Parse("https://example.com:8080/path/to/page?foo=bar&baz=123#section")
	if err != nil {
		return 0
	}
	if u.Scheme == "https" && u.Host == "example.com:8080" && u.Path == "/path/to/page" {
		return 1
	}
	return 0
}

// NetURLQuery tests query string parsing.
func NetURLQuery() int {
	u, _ := url.Parse("https://example.com?a=1&b=2&c=3")
	q := u.Query()
	if q.Get("a") == "1" && q.Get("b") == "2" && q.Get("c") == "3" {
		return 1
	}
	return 0
}

// NetURLQueryEncode tests Query().Encode().
func NetURLQueryEncode() int {
	q := url.Values{}
	q.Set("name", "hello world")
	q.Set("age", "30")
	encoded := q.Encode()
	if encoded != "" && len(encoded) > 0 {
		return 1
	}
	return 0
}

// NetURLResolveReference tests resolving relative URLs.
func NetURLResolveReference() int {
	base, _ := url.Parse("https://example.com/a/b/c")
	rel, _ := url.Parse("/d/e/f")
	resolved := base.ResolveReference(rel)
	if resolved.String() == "https://example.com/d/e/f" {
		return 1
	}
	return 0
}

// NetURLResolveReferenceRelative tests relative path resolution.
func NetURLResolveReferenceRelative() int {
	base, _ := url.Parse("https://example.com/a/b/")
	rel, _ := url.Parse("c/d")
	resolved := base.ResolveReference(rel)
	if resolved.Path == "/a/b/c/d" {
		return 1
	}
	return 0
}

// NetURLEscape tests URL path escaping.
func NetURLEscape() int {
	escaped := url.PathEscape("/path/with spaces/and?special=chars&")
	if escaped != "" {
		return 1
	}
	return 0
}

// NetURLUser tests URL with username and password.
func NetURLUser() int {
	u, _ := url.Parse("https://user:pass@example.com/")
	if u.User.Username() == "user" {
		pass, _ := u.User.Password()
		if pass == "pass" {
			return 1
		}
	}
	return 0
}

// NetURLString tests URL.String() reconstruction.
func NetURLString() int {
	u, _ := url.Parse("https://example.com:8080/path")
	if u.String() == "https://example.com:8080/path" {
		return 1
	}
	return 0
}

// ============================================================================
// net/netip — IP address parsing
// ============================================================================

// NetNetipParseAddr tests IPv4 address parsing.
func NetNetipParseAddr() int {
	addr, err := netip.ParseAddr("192.168.1.1")
	if err != nil || !addr.Is4() {
		return 0
	}
	return 1
}

// NetNetipParsePrefix tests IPv4 prefix parsing.
func NetNetipParsePrefix() int {
	pfx, err := netip.ParsePrefix("192.168.0.0/24")
	if err != nil || pfx.Bits() != 24 {
		return 0
	}
	return 1
}

// NetNetipIPv6 tests IPv6 address parsing.
func NetNetipIPv6() int {
	addr, err := netip.ParseAddr("::1")
	if err != nil || !addr.Is6() || !addr.IsLoopback() {
		return 0
	}
	return 1
}

// NetNetipIPv6Full tests full IPv6 address.
func NetNetipIPv6Full() int {
	addr, err := netip.ParseAddr("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	if err != nil || !addr.Is6() {
		return 0
	}
	return 1
}

// NetNetipFrom4 tests creating IPv4 from 4 bytes.
func NetNetipFrom4() int {
	addr := netip.AddrFrom4([4]byte{127, 0, 0, 1})
	if addr.IsLoopback() {
		return 1
	}
	return 0
}

// NetNetipCompare tests IP address comparison.
func NetNetipCompare() int {
	a1, _ := netip.ParseAddr("10.0.0.1")
	a2, _ := netip.ParseAddr("10.0.0.2")
	if a1.Less(a2) {
		return 1
	}
	return 0
}

// NetNetipPrefixContains tests prefix contains.
func NetNetipPrefixContains() int {
	pfx, _ := netip.ParsePrefix("192.168.0.0/24")
	ip1, _ := netip.ParseAddr("192.168.0.1")
	ip2, _ := netip.ParseAddr("192.168.1.1")
	if pfx.Contains(ip1) && !pfx.Contains(ip2) {
		return 1
	}
	return 0
}

// NetNetipMask tests IP prefix operations.
func NetNetipMask() int {
	addr, _ := netip.ParseAddr("192.168.1.1")
	prefix, _ := netip.ParsePrefix("192.168.0.0/16")
	if prefix.Contains(addr) {
		return 1
	}
	return 0
}
