package thirdparty

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
)

// ============================================================================
// crypto/md5
// ============================================================================

// CryptoMD5Sum tests MD5 hash of a string.
func CryptoMD5Sum() int {
	data := []byte("hello world")
	sum := md5.Sum(data)
	if len(sum) == 16 {
		return 1
	}
	return 0
}

// CryptoMD5Write tests md5.New() and Write.
func CryptoMD5Write() int {
	h := md5.New()
	h.Write([]byte("hello"))
	h.Write([]byte(" world"))
	sum := h.Sum(nil)
	if len(sum) == 16 {
		return 1
	}
	return 0
}

// ============================================================================
// crypto/sha1
// ============================================================================

// CryptoSHA1Sum tests SHA1 hash.
func CryptoSHA1Sum() int {
	data := []byte("hello world")
	sum := sha1.Sum(data)
	if len(sum) == 20 {
		return 1
	}
	return 0
}

// CryptoSHA1Write tests sha1.New() and Write.
func CryptoSHA1Write() int {
	h := sha1.New()
	h.Write([]byte("hello"))
	h.Write([]byte(" world"))
	sum := h.Sum(nil)
	if len(sum) == 20 {
		return 1
	}
	return 0
}

// ============================================================================
// crypto/sha256
// ============================================================================

// CryptoSHA256Sum tests SHA256 hash.
func CryptoSHA256Sum() int {
	data := []byte("hello world")
	sum := sha256.Sum256(data)
	if len(sum) == 32 {
		return 1
	}
	return 0
}

// CryptoSHA256Write tests sha256.New() and Write.
func CryptoSHA256Write() int {
	h := sha256.New()
	h.Write([]byte("hello"))
	h.Write([]byte(" world"))
	sum := h.Sum(nil)
	if len(sum) == 32 {
		return 1
	}
	return 0
}

// ============================================================================
// crypto/sha512
// ============================================================================

// CryptoSHA512Sum tests SHA512 hash.
func CryptoSHA512Sum() int {
	data := []byte("hello")
	sum := sha512.Sum512(data)
	if len(sum) == 64 {
		return 1
	}
	return 0
}

// CryptoSHA512_256 tests SHA512/256 (truncated).
func CryptoSHA512_256() int {
	data := []byte("hello")
	sum := sha512.Sum512_256(data)
	if len(sum) == 32 {
		return 1
	}
	return 0
}

// ============================================================================
// crypto/hmac
// ============================================================================

// CryptoHMAC_SHA256 tests HMAC-SHA256.
func CryptoHMAC_SHA256() int {
	key := []byte("secret key")
	message := []byte("hello world")
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	sum := mac.Sum(nil)
	if len(sum) == 32 {
		return 1
	}
	return 0
}

// CryptoHMACEqual tests HMAC comparison (constant-time).
func CryptoHMACEqual() int {
	key := []byte("secret key")
	message := []byte("hello")
	mac1 := hmac.New(sha256.New, key)
	mac1.Write(message)
	h1 := mac1.Sum(nil)

	mac2 := hmac.New(sha256.New, key)
	mac2.Write(message)
	h2 := mac2.Sum(nil)

	if hmac.Equal(h1, h2) {
		return 1
	}
	return 0
}

// CryptoHMACDifferent tests HMAC.Equal with wrong MAC.
func CryptoHMACDifferent() int {
	key := []byte("secret key")
	mac1 := hmac.New(sha256.New, key)
	mac1.Write([]byte("hello"))
	h1 := mac1.Sum(nil)

	h2 := make([]byte, 32)
	if hmac.Equal(h1, h2) {
		return 0
	}
	return 1
}

// ============================================================================
// crypto/aes + crypto/cipher
// ============================================================================

// CryptoAESEncrypt tests AES-CTR encryption and decryption.
func CryptoAESEncrypt() int {
	key := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	plaintext := []byte("hello world 1234")

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return 0
	}

	// CTR mode: encrypt
	ciphertext := make([]byte, len(plaintext))
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, plaintext)

	// CTR mode: decrypt (same operation)
	decrypted := make([]byte, len(ciphertext))
	stream2 := cipher.NewCTR(block, iv)
	stream2.XORKeyStream(decrypted, ciphertext)

	if bytes.Equal(decrypted, plaintext) {
		return 1
	}
	return 0
}

// CryptoAESCBC tests AES-CBC encryption and decryption.
func CryptoAESCBC() int {
	key := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	plaintext := make([]byte, 32)
	copy(plaintext, []byte("hello world 1234567890abcdef"))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return 0
	}

	// CBC mode encrypt
	blockSize := aes.BlockSize
	padded := make([]byte, blockSize)
	copy(padded, plaintext[:blockSize])
	ciphertext := make([]byte, blockSize)
	iv := make([]byte, blockSize)
	enc := cipher.NewCBCEncrypter(block, iv)
	enc.CryptBlocks(ciphertext, padded)

	// Decrypt
	decrypted := make([]byte, blockSize)
	dec := cipher.NewCBCDecrypter(block, ciphertext)
	dec.CryptBlocks(decrypted, ciphertext)

	if bytes.Equal(decrypted, padded) {
		return 1
	}
	return 0
}

// CryptoAESOFB tests AES-OFB mode.
func CryptoAESOFB() int {
	key := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	plaintext := []byte("hello world")
	block, _ := aes.NewCipher(key[:])
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewOFB(block, iv)

	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	// Decrypt
	stream2 := cipher.NewOFB(block, iv)
	decrypted := make([]byte, len(plaintext))
	stream2.XORKeyStream(decrypted, ciphertext)

	if bytes.Equal(decrypted, plaintext) {
		return 1
	}
	return 0
}

// ============================================================================
// crypto/subtle
// ============================================================================

// CryptoSubtleConstantTimeCompare tests constant-time comparison.
func CryptoSubtleConstantTimeCompare() int {
	a := []byte("hello world")
	b := []byte("hello world")
	c := []byte("hello worlc")
	if subtle.ConstantTimeCompare(a, b) == 1 && subtle.ConstantTimeCompare(a, c) == 0 {
		return 1
	}
	return 0
}

// CryptoSubtleConstantTimeCopy tests constant-time copy.
func CryptoSubtleConstantTimeCopy() int {
	src := []byte("hello world")
	dst := make([]byte, len(src))
	subtle.ConstantTimeCopy(5, dst, src)
	if string(dst[:5]) == "hello" {
		return 1
	}
	return 0
}

// CryptoSubtleXORBytes tests XOR of two slices.
func CryptoSubtleXORBytes() int {
	a := []byte{0xFF, 0x00, 0xAA, 0x55}
	b := []byte{0xAA, 0x55, 0xFF, 0x00}
	r := subtle.XORBytes(make([]byte, 4), a, b)
	// r == 4 means success (4 bytes processed)
	if r == 4 {
		return 1
	}
	return 0
}
