package thirdparty

import (
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
)

// HashAdler32 tests Adler-32 checksum.
func HashAdler32() int {
	data := []byte("hello world")
	checksum := adler32.Checksum(data)
	if checksum != 0 {
		return 1
	}
	return 0
}

// HashAdler32Write tests New and Write.
func HashAdler32Write() int {
	h := adler32.New()
	h.Write([]byte("hello"))
	h.Write([]byte(" world"))
	sum := h.Sum32()
	if sum != 0 {
		return 1
	}
	return 0
}

// HashCrc32 tests CRC32-IEEE checksum.
func HashCrc32() int {
	data := []byte("hello world")
	h := crc32.NewIEEE()
	h.Write(data)
	sum := h.Sum32()
	if sum != 0 {
		return 1
	}
	return 0
}

// HashCrc32IEEE tests CRC32 with IEEE polynomial.
func HashCrc32IEEE() int {
	data := []byte("hello")
	checksum := crc32.ChecksumIEEE(data)
	if checksum != 0 {
		return 1
	}
	return 0
}

// HashCrc64ECMA tests CRC64 with ECMA polynomial.
func HashCrc64ECMA() int {
	data := []byte("hello")
	h := crc64.New(crc64.MakeTable(crc64.ECMA))
	h.Write(data)
	sum := h.Sum64()
	if sum != 0 {
		return 1
	}
	return 0
}

// HashCrc64ISO tests CRC64 with ISO polynomial.
func HashCrc64ISO() int {
	data := []byte("test")
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	h.Write(data)
	sum := h.Sum64()
	if sum != 0 {
		return 1
	}
	return 0
}

// HashFnv64 tests FNV-1 64-bit hash.
func HashFnv64() int {
	h := fnv.New64()
	h.Write([]byte("hello"))
	sum := h.Sum64()
	if sum != 0 {
		return 1
	}
	return 0
}

// HashFnv64a tests FNV-1a 64-bit hash.
func HashFnv64a() int {
	h := fnv.New64a()
	h.Write([]byte("hello"))
	sum := h.Sum64()
	if sum != 0 {
		return 1
	}
	return 0
}
