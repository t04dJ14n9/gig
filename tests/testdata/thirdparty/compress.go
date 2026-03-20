package thirdparty

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"io"
)


// BinaryWriteRead tests encoding/binary Write and Read.
func BinaryWriteRead() int {
	var buf bytes.Buffer
	var val uint32 = 0x12345678
	err := binary.Write(&buf, binary.BigEndian, val)
	if err != nil {
		return 0
	}
	var result uint32
	err = binary.Read(&buf, binary.BigEndian, &result)
	if err != nil {
		return 0
	}
	if result == 0x12345678 {
		return 1
	}
	return 0
}

// BinaryPutGet tests PutUint32 and GetUint32.
func BinaryPutGet() int {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, 0xDEADBEEF)
	val := binary.BigEndian.Uint32(data)
	if val == 0xDEADBEEF {
		return 1
	}
	return 0
}

// BinaryMultipleValues tests writing multiple values.
func BinaryMultipleValues() int {
	var buf bytes.Buffer
	var a uint16 = 0x1234
	var b uint32 = 0x56789ABC
	var c uint64 = 0xDEF0123456789ABC
	binary.Write(&buf, binary.BigEndian, a)
	binary.Write(&buf, binary.BigEndian, b)
	binary.Write(&buf, binary.BigEndian, c)

	var ra uint16
	var rb4 uint32
	var rc8 uint64
	binary.Read(&buf, binary.BigEndian, &ra)
	binary.Read(&buf, binary.BigEndian, &rb4)
	binary.Read(&buf, binary.BigEndian, &rc8)
	if ra == a && rb4 == b && rc8 == c {
		return 1
	}
	return 0
}

// BinaryLittleEndian tests little-endian byte order.
func BinaryLittleEndian() int {
	var buf bytes.Buffer
	val := uint32(0x12345678)
	binary.Write(&buf, binary.LittleEndian, val)

	var result uint32
	binary.Read(&buf, binary.LittleEndian, &result)
	if result == 0x12345678 {
		return 1
	}
	return 0
}

// CompressGzipRoundtrip tests gzip compress then decompress.
func CompressGzipRoundtrip() int {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	w.Write([]byte("hello gzip"))
	w.Close()

	r, err := gzip.NewReader(&buf)
	if err != nil {
		return 0
	}
	data, err := io.ReadAll(r)
	r.Close()
	if err != nil {
		return 0
	}
	if string(data) == "hello gzip" {
		return 1
	}
	return 0
}

// CompressZlibRoundtrip tests zlib compress then decompress.
func CompressZlibRoundtrip() int {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write([]byte("hello zlib"))
	w.Close()

	r, err := zlib.NewReader(&buf)
	if err != nil {
		return 0
	}
	data, err := io.ReadAll(r)
	r.Close()
	if err != nil {
		return 0
	}
	if string(data) == "hello zlib" {
		return 1
	}
	return 0
}

