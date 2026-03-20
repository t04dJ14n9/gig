package thirdparty

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
)

// JsonMarshal tests json.Marshal.
func JsonMarshal() int {
	data := map[string]int{"a": 1, "b": 2}
	b, _ := json.Marshal(data)
	return len(b)
}

// JsonUnmarshal tests json.Unmarshal.
func JsonUnmarshal() int {
	var data map[string]int
	json.Unmarshal([]byte(`{"a":1,"b":2}`), &data)
	return data["a"] + data["b"]
}

// JsonMarshalIndent tests json.MarshalIndent.
func JsonMarshalIndent() int {
	data := map[string]int{"a": 1}
	b, _ := json.MarshalIndent(data, "", "  ")
	return len(b)
}

// JsonDecode tests json.NewDecoder.
func JsonDecode() int {
	var data map[string]int
	decoder := json.NewDecoder(strings.NewReader(`{"x":10}`))
	decoder.Decode(&data)
	return data["x"]
}

// JsonEncode tests json.NewEncoder.
func JsonEncode() int {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(map[string]int{"y": 20})
	return buf.Len()
}

// JsonNumber tests json.Number.
func JsonNumber() int {
	var data map[string]json.Number
	json.Unmarshal([]byte(`{"n":"42"}`), &data)
	n, _ := data["n"].Int64()
	return int(n)
}

// Base64Encode tests base64.StdEncoding.EncodeToString.
func Base64Encode() string {
	return base64.StdEncoding.EncodeToString([]byte("hello"))
}

// Base64Decode tests base64.StdEncoding.DecodeString.
func Base64Decode() string {
	decoded, _ := base64.StdEncoding.DecodeString("aGVsbG8=")
	return string(decoded)
}

// Base64URLEncode tests base64.URLEncoding.EncodeToString.
func Base64URLEncode() string {
	return base64.URLEncoding.EncodeToString([]byte("hello world"))
}

// Base64URLDecode tests base64.URLEncoding.DecodeString.
func Base64URLDecode() string {
	decoded, _ := base64.URLEncoding.DecodeString("aGVsbG8gd29ybGQ=")
	return string(decoded)
}

// Base64NewEncoder tests base64.NewEncoder.
func Base64NewEncoder() int {
	buf := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	encoder.Write([]byte("test"))
	encoder.Close()
	return buf.Len()
}

// Base64NewDecoder tests base64.NewDecoder.
func Base64NewDecoder() int {
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader("aGVsbG8="))
	result, _ := io.ReadAll(decoder)
	return len(result)
}

// HexEncodeToString tests hex.EncodeToString.
func HexEncodeToString() string {
	return hex.EncodeToString([]byte("hello"))
}

// HexDecodeString tests hex.DecodeString.
func HexDecodeString() string {
	decoded, _ := hex.DecodeString("68656c6c6f")
	return string(decoded)
}

// HexNewEncoder tests hex.NewEncoder.
func HexNewEncoder() int {
	buf := new(bytes.Buffer)
	encoder := hex.NewEncoder(buf)
	encoder.Write([]byte("test"))
	return buf.Len()
}

// HexNewDecoder tests hex.NewDecoder.
func HexNewDecoder() int {
	decoder := hex.NewDecoder(strings.NewReader("74657374"))
	result, _ := io.ReadAll(decoder)
	return len(result)
}
