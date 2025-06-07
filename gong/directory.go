package gong

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"github.com/asciifaceman/gong/gong/magic"
	"github.com/spf13/afero"
)

// LoadGongDirectory reads a Gong file and returns the directory entries
// as a slice of Asset pointers, the total size of the directory, and any error encountered.
func LoadGongDirectory(path string) ([]*Asset, uint32, error) {
	f := afero.NewOsFs()
	file, err := f.Open(path)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to open Gong file %s: %w", path, err)
	}
	defer file.Close()

	// verify SOG + IDENT
	header := make([]byte, len(magic.SOG)+len(magic.GONG))
	if _, err := file.Read(header); err != nil {
		return nil, 0, fmt.Errorf("failed to read Gong file: %w", err)
	}

	if !reflect.DeepEqual(header[:len(magic.SOG)], magic.SOG) || !reflect.DeepEqual(header[len(magic.SOG):len(magic.SOG)+len(magic.GONG)], magic.GONG) {
		return nil, 0, fmt.Errorf("invalid Gong file header: %s", path)
	}

	// read dircount and dirsize
	h := magic.GongHeader{}
	if err := binary.Read(file, binary.BigEndian, &h); err != nil {
		return nil, 0, fmt.Errorf("failed to read Gong header: %w", err)
	}

	// parse directory entries
	entries := make([]*Asset, 0, h.DirCount)
	startOffset := int64(magic.HeaderSize)
	if _, err := file.Seek(startOffset, io.SeekStart); err != nil {
		return nil, 0, fmt.Errorf("failed to seek to directory entries: %w", err)
	}

	var maxAssetEnd uint32 = 0

	for i := 0; i < int(h.DirCount); i++ {
		idLen := make([]byte, 1)
		if _, err := file.Read(idLen); err != nil {
			return nil, 0, fmt.Errorf("failed to read ID length for entry %d: %w", i, err)
		}
		idBytes := make([]byte, idLen[0])
		if _, err := file.Read(idBytes); err != nil {
			return nil, 0, fmt.Errorf("failed to read ID for entry %d: %w", i, err)
		}

		nameLen := make([]byte, 1)
		if _, err := file.Read(nameLen); err != nil {
			return nil, 0, fmt.Errorf("failed to read name length for entry %d: %w", i, err)
		}

		nameBytes := make([]byte, nameLen[0])
		if _, err := file.Read(nameBytes); err != nil {
			return nil, 0, fmt.Errorf("failed to read name for entry %d: %w", i, err)
		}

		var offset uint32
		var size uint32
		var compression byte

		if err := binary.Read(file, binary.BigEndian, &offset); err != nil {
			return nil, 0, fmt.Errorf("failed to read offset for entry %d: %w", i, err)
		}
		if err := binary.Read(file, binary.BigEndian, &size); err != nil {
			return nil, 0, fmt.Errorf("failed to read size for entry %d: %w", i, err)
		}
		if err := binary.Read(file, binary.BigEndian, &compression); err != nil {
			return nil, 0, fmt.Errorf("failed to read compression for entry %d: %w", i, err)
		}

		entries = append(entries, &Asset{
			ID:          string(idBytes),
			Filename:    string(nameBytes),
			Offset:      offset,
			Size:        size,
			Compression: compression,
		})

		end := offset + size
		if end > maxAssetEnd {
			maxAssetEnd = end
		}
	}

	return entries, h.DirSize, nil
}

// CalculateDirSize computes the total size of the directory entries
// based on the assets provided. It returns the size in bytes.
func CalculateDirSize(assets []*Asset) uint32 {
	var size uint32

	for _, asset := range assets {
		if asset == nil {
			continue
		}
		size += 1 + uint32(len(asset.ID)) +
			1 + uint32(len(asset.Filename)) +
			4 + 4 + 1 // offset, size, and compression type
	}

	return size
}
