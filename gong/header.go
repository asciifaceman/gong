package gong

import (
	"encoding/binary"
	"io"

	"github.com/spf13/afero"
)

var (
	Endianess = binary.BigEndian // bigendian for encoding

	NULL byte = 0x00 // NULL byte for string termination

	SOG  Marker = []byte{0xFF, 0xD6}                   // Start of Gong
	EOG  Marker = []byte{0xFF, 0xD7}                   // End of Gong
	BOM  Marker = []byte{0xFE, 0xFF}                   // Byte Order Mark for BigEndian
	GONG Marker = []byte{0x47, 0x4F, 0x4E, 0x47, 0x00} // Gong file marker

	SODE Marker = []byte{0xFF, 0xD8} // Start of Directory Entry
)

// Marker is a binary marker used to identify files/sections
type Marker []byte

func (m Marker) Is(magic []byte) bool {
	if len(m) != len(magic) {
		return false
	}
	for i := range magic {
		if m[i] != magic[i] {
			return false
		}
	}
	return true
}

// IsAt checks if the header at the specified offset in the file matches the provided magic bytes.
func (m Marker) IsAt(f afero.File, offset int64) bool {
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return false
	}

	header := make([]byte, len(m))
	if _, err := f.Read(header); err != nil {
		return false
	}
	if !m.Is(header) {
		return false
	}
	return true
}

type GongHeader struct {
	SOG      [2]byte // Start of Gong marker
	BOM      [2]byte // Byte Order Mark/Endianess
	GONG     [5]byte // Gong file marker
	VER      [2]byte // Version of the gong file format
	RESERVED [2]byte // Reserved for future use 0x00
	COUNT    uint32  // Number of directory entries
	OFFSET   uint32  // Offset in the gong file where the directory starts
}
