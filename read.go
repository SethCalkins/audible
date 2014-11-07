// This file is subject to a BSD license.
// Its contents can be found in the enclosed LICENSE file.

package audible

import (
	"encoding/binary"
	"io"
)

var endian = binary.BigEndian

// readU8 reads a uint8 value.
func readU8(r io.Reader) (v uint8, err error) {
	err = binary.Read(r, endian, &v)
	return v, err
}

// readU32 reads a uint32 value.
func readU32(r io.Reader) (v uint32, err error) {
	err = binary.Read(r, endian, &v)
	return v, err
}

// readString reads a string value of the given length
func readString(r io.Reader, size uint32) (string, error) {
	v, err := readBytes(r, size)
	return string(v), err
}

// readBytes reads a byte slice of the given length.
func readBytes(r io.Reader, size uint32) ([]byte, error) {
	buf := make([]byte, size)
	_, err := io.ReadFull(r, buf)
	return buf, err
}
