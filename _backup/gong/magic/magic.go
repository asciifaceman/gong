package magic

import (
	"encoding/binary"

	"github.com/spf13/afero"
)

var (
	NULL byte = 0x00                                 // NULL byte
	SOG       = []byte{0xFF, 0xD6}                   // SOG (Start of Gong) magic number
	EOG       = []byte{0xff, 0xD7}                   // EOG (End of Gong) magic number
	GONG      = []byte{0x47, 0x4F, 0x4E, 0x47, 0x00} // GONG magic number
	SOA       = []byte{0xFF, 0xD8}                   // SOA (Start of Asset) magic number
	EOA       = []byte{0xFF, 0xD9}                   // EOA (End of Asset) magic number
	ASST      = []byte{0x41, 0x53, 0x53, 0x54}       // ASST (Assset Leading Header) magic number

	HeaderSize uint32 = uint32(len(SOG) + len(GONG) + 2 + 4) // SOG + GONG + DirCount (2 bytes) + DirSize (4 bytes
)

type GongHeader struct {
	DirCount uint16
	DirSize  uint32
}

func WriteLeadingMagic(f afero.File) error {
	if _, err := f.Write(SOG); err != nil {
		return err
	}
	if _, err := f.Write(GONG); err != nil {
		return err
	}
	header := GongHeader{
		DirCount: 0,
		DirSize:  0,
	}

	if err := binary.Write(f, binary.BigEndian, header); err != nil {
		return err
	}
	return nil
}

func WriteTailingMagic(f afero.File) error {
	if _, err := f.Write(EOG); err != nil {
		return err
	}
	return nil
}
