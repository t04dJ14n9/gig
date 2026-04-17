package divergence_hunt113

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// ============================================================================
// Round 113: Encoding/hex, base64 roundtrips
// ============================================================================

func HexEncode() string {
	data := []byte("hello")
	return hex.EncodeToString(data)
}

func HexDecode() string {
	data, _ := hex.DecodeString("68656c6c6f")
	return string(data)
}

func HexRoundtrip() string {
	original := []byte("test data")
	encoded := hex.EncodeToString(original)
	decoded, _ := hex.DecodeString(encoded)
	return fmt.Sprintf("%v", string(decoded) == string(original))
}

func Base64Encode() string {
	data := []byte("hello world")
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode() string {
	data, _ := base64.StdEncoding.DecodeString("aGVsbG8gd29ybGQ=")
	return string(data)
}

func Base64Roundtrip() string {
	original := []byte("test data for base64")
	encoded := base64.StdEncoding.EncodeToString(original)
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	return fmt.Sprintf("%v", string(decoded) == string(original))
}

func Base64URLEncoding() string {
	data := []byte("url safe data")
	encoded := base64.URLEncoding.EncodeToString(data)
	decoded, _ := base64.URLEncoding.DecodeString(encoded)
	return fmt.Sprintf("%v", string(decoded) == string(data))
}

func HexEmpty() string {
	return hex.EncodeToString([]byte{})
}

func Base64Empty() string {
	return base64.StdEncoding.EncodeToString([]byte{})
}

func HexEncodeNumbers() string {
	data := []byte{0x00, 0xff, 0x0f, 0xf0}
	return hex.EncodeToString(data)
}
