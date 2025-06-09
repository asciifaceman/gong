package gong

import "encoding/binary"

var (
	Endianess = binary.BigEndian // bigendian for encoding

	NULL byte = 0x00 // NULL byte for string termination

	SOG  Marker = []byte{0xFF, 0xD6}                   // Start of Gong
	EOG  Marker = []byte{0xFF, 0xD7}                   // End of Gong
	BOM  Marker = []byte{0xFE, 0xFF}                   // Byte Order Mark for BigEndian
	GONG Marker = []byte{0x47, 0x4F, 0x4E, 0x47, 0x00} // Gong file marker

)

// Marker is a binary marker used to identify files/sections
type Marker []byte

type GongHeader struct {
	SOG      [2]byte // Start of Gong marker
	BOM      [2]byte // Byte Order Mark/Endianess
	GONG     [5]byte // Gong file marker
	VER      [2]byte // Version of the gong file format
	RESERVED [2]byte // Reserved for future use 0x00
	COUNT    uint32  // Number of directory entries
	OFFSET   uint32  // Offset in the gong file where the directory starts
}

type DirectoryEntryHeader struct {
	SODE        [2]byte // Start of Directory entry marker
	SIZE        uint32  // Size of the directory entry header
	ID_LEN      uint16  // Length of the ID string
	ID          []byte  // ID of the asset (unique identifier)
	FNAME_LEN   uint16  // Length of the filename string
	FNAME       []byte  // Filename of the asset sans path and extension
	FileTypeLen uint16  // Length of the file type string
	FileType    []byte  // File type of the asset (extension)
	Offset      uint32  // Offset in the gong file where the content starts
	Size        uint32  // Size of the content stored at offset
	Compression byte    // Compression type (0 = none)
}
