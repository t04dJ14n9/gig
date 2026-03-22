package thirdparty

import (
	"mime"
	"mime/multipart"
	"strings"
)

// ============================================================================
// mime — MIME type detection
// ============================================================================

// MimeTypeByExtension tests MIME type lookup by extension.
func MimeTypeByExtension() int {
	t := mime.TypeByExtension(".html")
	if t == "text/html; charset=utf-8" || t == "text/html" {
		return 1
	}
	return 0
}

// MimeTypeByExtensionJSON tests MIME type for JSON.
func MimeTypeByExtensionJSON() int {
	t := mime.TypeByExtension(".json")
	if t == "application/json" {
		return 1
	}
	return 0
}

// MimeExtensionsByType tests extension lookup by MIME type.
func MimeExtensionsByType() int {
	exts, _ := mime.ExtensionsByType("text/html")
	if len(exts) > 0 && exts[0] == ".html" {
		return 1
	}
	return 0
}

// MimeParseMediaType tests parsing MIME type strings.
func MimeParseMediaType() int {
	mediaType, params, err := mime.ParseMediaType("text/html; charset=utf-8")
	if err != nil || mediaType != "text/html" || params["charset"] != "utf-8" {
		return 0
	}
	return 1
}

// MimeWordEncoder tests MIME word encoding (RFC 2047).
func MimeWordEncoder() int {
	encoded := mime.QEncoding.Encode("utf-8", "hello world")
	if strings.Contains(encoded, "?utf-8?") {
		return 1
	}
	return 0
}

// MimeWordDecoder tests MIME word decoding via WordDecoder.
func MimeWordDecoder() int {
	dec := new(mime.WordDecoder)
	decoded, err := dec.Decode("=?utf-8?q?=48=65=6c=6c=6f?=")
	if err != nil || decoded != "Hello" {
		return 0
	}
	return 1
}

// ============================================================================
// mime/multipart — multipart form parsing
// ============================================================================

// MultipartNewReader tests multipart form reading.
func MultipartNewReader() int {
	body := "--boundary\r\nContent-Disposition: form-data; name=\"field\"\r\n\r\nvalue\r\n--boundary--\r\n"
	r := multipart.NewReader(strings.NewReader(body), "boundary")
	part, err := r.NextPart()
	if err != nil || part.FormName() != "field" {
		return 0
	}
	data := make([]byte, 100)
	n, _ := part.Read(data)
	if string(data[:n]) == "value" {
		return 1
	}
	return 0
}

// MultipartCreate tests creating multipart form.
func MultipartCreate() int {
	var buf strings.Builder
	w := multipart.NewWriter(&buf)
	w.SetBoundary("myboundary")

	// Write a text field
	part, _ := w.CreateFormField("name")
	part.Write([]byte("Alice"))

	// Write a file
	filePart, _ := w.CreateFormFile("file", "test.txt")
	filePart.Write([]byte("hello world"))

	w.Close()

	content := buf.String()
	if strings.Contains(content, "myboundary") && strings.Contains(content, "Alice") {
		return 1
	}
	return 0
}

// MultipartFormData tests form-data content disposition.
func MultipartFormData() int {
	body := "--boundary\r\nContent-Disposition: form-data; name=\"file\"; filename=\"test.txt\"\r\n\r\ndata\r\n--boundary--\r\n"
	r := multipart.NewReader(strings.NewReader(body), "boundary")
	part, _ := r.NextPart()
	if part.FileName() == "test.txt" && part.FormName() == "file" {
		return 1
	}
	return 0
}
