package gong

import "bytes"

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

func (d *DirectoryEntryHeader) EncodedSize() uint32 {
	size := 2 + // SODE
		4 + // SIZE
		2 + // ID_LEN
		uint32(len(d.ID)) + // ID
		2 + // FNAME_LEN
		uint32(len(d.FNAME)) + // FNAME
		2 + // FileTypeLen
		uint32(len(d.FileType)) + // FileType
		4 + // Offset
		4 + // Size
		1 // Compression

	return size
}

func (d *DirectoryEntryHeader) Encode() []byte {
	var enc bytes.Buffer

}
