package gong

import (
	"bytes"
	"encoding/binary"
)

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
		1 + // Compression
		2 // EODE (End of Directory entry marker)

	return size
}

func (d *DirectoryEntryHeader) Encode() []byte {
	var enc bytes.Buffer

	size := d.EncodedSize()
	enc.Grow(int(size))

	if _, err := enc.Write(SODE); err != nil {
		return nil // Error writing SODE, return empty slice
	}

	// write size uint32 bytes
	if err := binary.Write(&enc, Endianess, d.EncodedSize()); err != nil {
		return nil // Error writing SIZE, return empty slice
	}

	// write ID length uint16 and ID []byte
	idLen := len(d.ID)
	if err := binary.Write(&enc, Endianess, uint16(idLen)); err != nil {
		return nil // Error writing ID length, return empty slice
	}
	if _, err := enc.Write(d.ID); err != nil {
		return nil // Error writing ID bytes, return empty slice
	}

	// write FNAME length uint16 and FNAME []byte
	fnameLen := len(d.FNAME)
	if err := binary.Write(&enc, Endianess, uint16(fnameLen)); err != nil {
		return nil // Error writing FNAME length, return empty slice
	}
	if _, err := enc.Write(d.FNAME); err != nil {
		return nil // Error writing FNAME bytes, return empty slice
	}

	// write FileType length uint16 and FileType []byte
	fileTypeLen := len(d.FileType)
	if err := binary.Write(&enc, Endianess, uint16(fileTypeLen)); err != nil {
		return nil // Error writing FileType length, return empty slice
	}
	if _, err := enc.Write(d.FileType); err != nil {
		return nil // Error writing FileType bytes, return empty slice
	}

	// write Offset uint32
	if err := binary.Write(&enc, Endianess, d.Offset); err != nil {
		return nil // Error writing Offset, return empty slice
	}

	// write Size uint32
	if err := binary.Write(&enc, Endianess, d.Size); err != nil {
		return nil // Error writing Size, return empty slice
	}

	// write Compression byte
	if err := enc.WriteByte(d.Compression); err != nil {
		return nil // Error writing Compression, return empty slice
	}

	// write EODE marker
	if _, err := enc.Write(EODE); err != nil {
		return nil // Error writing EODE, return empty slice
	}

	return enc.Bytes()
}
