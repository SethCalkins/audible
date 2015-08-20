// This file is subject to a BSD license.
// Its contents can be found in the enclosed LICENSE file.

package audible

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

// fileMagic defines the magic number that identifies an audible file.
const fileMagic = 1469084982

// Header defines the header for an Audible .aa file.
type Header struct {
	Filesize   uint32            // Total file size.
	Magic      uint32            // File magic value.
	TOC        [][2]uint32       // Table of contents defines the offsets and size of various blocks.
	Tags       map[string]string // Table of key/value tag pairs.
	HeaderSeed uint32
	HeaderKey  []byte
}

// ReadFile reads metadata for the given file.
func ReadFile(file string) (*Header, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return Read(fd)
}

// Read reads metadata for the given reader.
func Read(r io.Reader) (hdr *Header, err error) {
	hdr = new(Header)
	hdr.Tags = make(map[string]string)

	hdr.Filesize, err = readU32(r)
	if err != nil {
		return nil, err
	}

	// Ensure we have a valid audible file.
	hdr.Magic, err = readU32(r)
	if err != nil {
		return nil, err
	}

	if hdr.Magic != fileMagic {
		return nil, errors.New("not a valid audible file")
	}

	// Read the table of contents
	tocSize, err := readU32(r)
	if err != nil {
		return nil, err
	}
	hdr.TOC = make([][2]uint32, tocSize)

	// Unidentified integer.
	_, err = readU32(r)
	if err != nil {
		return nil, err
	}

	for i := range hdr.TOC {
		// TOC entry index
		_, err = readU32(r)
		if err != nil {
			return nil, err
		}

		// Block Offset.
		hdr.TOC[i][0], err = readU32(r)
		if err != nil {
			return nil, err
		}

		// Block size.
		hdr.TOC[i][1], err = readU32(r)
		if err != nil {
			return nil, err
		}
	}

	// Header termination block.
	_, err = readBytes(r, 24)
	if err != nil {
		return nil, err
	}

	// Read dictionary entries.
	npairs, err := readU32(r)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(npairs); i++ {
		// Unidentified byte.
		_, err = readU8(r)
		if err != nil {
			return nil, err
		}

		// Length of key string.
		nkey, err := readU32(r)
		if err != nil {
			return nil, err
		}

		// Length of value string.
		nval, err := readU32(r)
		if err != nil {
			return nil, err
		}

		// Key string.
		key, err := readString(r, nkey)
		if err != nil {
			return nil, err
		}

		// Value string.
		val, err := readString(r, nval)
		if err != nil {
			return nil, err
		}

		if key == "HeaderSeed" {
			i, _ := strconv.Atoi(val)
			hdr.HeaderSeed = (uint32)(i)
		}

		if key == "HeaderKey" {
			data := strings.Split(val, " ")
			buf := new(bytes.Buffer)
			for _, item := range data {
				i, _ := strconv.Atoi(item)
				_ = binary.Write(buf, binary.BigEndian, (uint32)(i))
			}
			hdr.HeaderKey = buf.Bytes()
		}

		hdr.Tags[key] = val
	}

	return hdr, nil
}
