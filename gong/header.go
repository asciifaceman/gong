package gong

import (
	"encoding/binary"
	"io"

	"github.com/spf13/afero"
)

var (
	Endianess = binary.BigEndian // Endianness used for gong file headers

	NULL byte   = 0x00                                 // NULL byte
	SOG  Header = []byte{0xFF, 0xD6}                   // SOG (Start of Gong) magic number, should be bytes 0 and 1
	EOG  Header = []byte{0xff, 0xD7}                   // EOG (End of Gong) magic number
	GONG Header = []byte{0x47, 0x4F, 0x4E, 0x47, 0x00} // GONG magic number, should be bytes 2-6
	SOA  Header = []byte{0xFF, 0xD8}                   // SOA (Start of Asset) magic number
	EOA  Header = []byte{0xFF, 0xD9}                   // EOA (End of Asset) magic number

	// DirectoryOffset is the offset in bytes where the directory entries start in a gong file.
	DirectoryOffset = int64(len(SOG) + len(GONG)) // Size of the preamble (SOG + GONG)

	// AssetStartOffset is the offset in bytes where the first asset starts in a gong file.
	AssetStartOffset = DirectoryOffset + DirectoryHeader{}.Size() // Size of the directory header

	// DefaultDirectoryHeader is the default directory header used when no directory entries are present.
	DefaultDirectoryHeader = DirectoryHeader{
		DirCount: 0, // Default directory count
		DirSize:  0, // Default directory size
	}
)

// Header represents a generic header in a gong file.
// It is a byte slice that contains the magic bytes identifying the file type.
type Header []byte

// Is checks if the header matches the provided magic bytes.
func (h Header) Is(magic []byte) bool {
	if len(h) != len(magic) {
		return false
	}
	for i := range magic {
		if h[i] != magic[i] {
			return false
		}
	}
	return true
}

// Size returns the size of the header in bytes.
func (h Header) Size() uint32 {
	return uint32(len(h))
}

// Write writes the header to the provided file.
func (h Header) Write(f afero.File) error {
	if _, err := f.Write(h); err != nil {
		return err
	}
	return nil
}

// WriteAt writes the header to the provided file at the specified offset.
func (h Header) WriteAt(f afero.File, offset int64) error {
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return err
	}
	return h.Write(f)
}

// IsAt checks if the header at the specified offset in the file matches the provided magic bytes.
func (h Header) IsAt(f afero.File, offset int64) bool {
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return false
	}

	header := make([]byte, len(h))
	if _, err := f.Read(header); err != nil {
		return false
	}
	if !h.Is(header) {
		return false
	}
	return true
}

func ReadDirectoryHeader(f afero.File) (DirectoryHeader, error) {
	var header DirectoryHeader
	if _, err := f.Seek(DirectoryOffset, io.SeekStart); err != nil {
		return header, err
	}

	err := binary.Read(f, Endianess, &header)
	return header, err
}

// DirectoryHeader represents the header for the directory entries in a gong file.
// It contains the number of directory entries and their total size.
// should be bytes 7-12 in the gong file.
type DirectoryHeader struct {
	DirCount uint16 // Number of directory entries
	DirSize  uint32 // Size of the directory entries in bytes
}

// Size returns the size of the DirectoryHeader in bytes.
func (h DirectoryHeader) Size() int64 {
	return 2 + 4 // 2 bytes for DirCount + 4 bytes for DirSize
}

func (h DirectoryHeader) Patch(f afero.File, dirCount uint16, dirSize uint32) error {
	h.DirCount = dirCount
	h.DirSize = dirSize
	if _, err := f.Seek(DirectoryOffset, io.SeekStart); err != nil { // Seek to the start of the DirectoryHeader
		return err
	}
	return h.Write(f)
}

// Write writes the DirectoryHeader to the provided file.
func (h DirectoryHeader) Write(f afero.File) error {
	return binary.Write(f, Endianess, h)
}

// WriteAt writes the DirectoryHeader to the provided file at the specified offset.
func (h DirectoryHeader) WriteAt(f afero.File, offset int64) error {
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return err
	}
	return h.Write(f)
}
