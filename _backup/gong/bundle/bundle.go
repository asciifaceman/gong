package bundle

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/asciifaceman/gong/gong/assets"
	"github.com/asciifaceman/gong/gong/gio"
	"github.com/asciifaceman/gong/gong/magic"
	"github.com/spf13/afero"
)

func CreateEmptyGong(path string) error {
	f := afero.NewOsFs()
	if gio.Exists(f, path) {
		return gio.ErrFileExists
	}

	// Create an empty file at the specified path
	file, err := f.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the magic header
	if err := magic.WriteLeadingMagic(file); err != nil {
		return err
	}

	header := magic.GongHeader{
		DirCount: 0,
		DirSize:  0,
	}

	if err := binary.Write(file, binary.BigEndian, header); err != nil {
		return err
	}

	if err := magic.WriteTailingMagic(file); err != nil {
		return err
	}
	return nil
}

func CalculateDirSize(assets []*assets.Asset) uint32 {
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

func WriteGongFile(path string, ast []*assets.Asset) error {
	f := afero.NewOsFs()

	var (
		file     afero.File
		existing = make([]*assets.Asset, 0) // existing assets in the Gong file
		assetDB  uint32
		err      error
	)

	if gio.Exists(f, path) {
		// load existing gong file
		existing, assetDB, err = LoadGongDirectory(path)
		if err != nil {
			return fmt.Errorf("failed to load existing Gong file %s: %w", path, err)
		}
		file, err = f.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("failed to open Gong file %s for writing: %w", path, err)
		}
	} else {
		file, err = f.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create Gong file %s: %w", path, err)
		}
		assetDB = magic.HeaderSize
		if err = magic.WriteLeadingMagic(file); err != nil {
			file.Close()
			return fmt.Errorf("failed to write leading magic to Gong file %s: %w", path, err)
		}
	}

	defer file.Close()

	allEntries := append(existing, ast...)

	for i := range ast {
		ast[i].Offset = assetDB
		assetDB += ast[i].Size
	}

	if _, err := file.Seek(int64(magic.HeaderSize), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek to Gong header in file %s: %w", path, err)
	}

	for _, entry := range allEntries {
		if len(entry.ID) > 255 || len(entry.Filename) > 255 {
			return fmt.Errorf("ID or filename too long for entry %s: ID length %d, Filename length %d", entry.ID, len(entry.ID), len(entry.Filename))
		}

		if _, err := file.Write([]byte{byte(len(entry.ID))}); err != nil {
			return fmt.Errorf("failed to write ID length for entry %s: %w", entry.ID, err)
		}
		if _, err := file.Write([]byte(entry.ID)); err != nil {
			return fmt.Errorf("failed to write ID for entry %s: %w", entry.ID, err)
		}
		if _, err := file.Write([]byte{byte(len(entry.Filename))}); err != nil {
			return fmt.Errorf("failed to write filename length for entry %s: %w", entry.ID, err)
		}
		if _, err := file.Write([]byte(entry.Filename)); err != nil {
			return fmt.Errorf("failed to write filename for entry %s: %w", entry.ID, err)
		}
		if err := binary.Write(file, binary.BigEndian, entry.Offset); err != nil {
			return fmt.Errorf("failed to write offset for entry %s: %w", entry.ID, err)
		}
		if err := binary.Write(file, binary.BigEndian, entry.Size); err != nil {
			return fmt.Errorf("failed to write size for entry %s: %w", entry.ID, err)
		}
		if _, err := file.Write([]byte{entry.Compression}); err != nil {
			return fmt.Errorf("failed to write compression for entry %s: %w", entry.ID, err)
		}
	}

	// patch header
	if _, err := file.Seek(int64(len(magic.SOG)+len(magic.GONG)), io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek to Gong header in file %s: %w", path, err)
	}
	dirSize := CalculateDirSize(allEntries)
	if err := binary.Write(file, binary.BigEndian, uint16(len(allEntries))); err != nil {
		return fmt.Errorf("failed to write directory count for Gong file %s: %w", path, err)
	}
	if err := binary.Write(file, binary.BigEndian, dirSize); err != nil {
		return fmt.Errorf("failed to write directory size for Gong file %s: %w", path, err)
	}

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("failed to seek to end of Gong file %s: %w", path, err)
	}
	if err := magic.WriteTailingMagic(file); err != nil {
		return fmt.Errorf("failed to write trailing magic to Gong file %s: %w", path, err)
	}
	return nil

}

func LoadGongDirectory(path string) ([]*assets.Asset, uint32, error) {
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
	entries := make([]*assets.Asset, 0, h.DirCount)
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

		entries = append(entries, &assets.Asset{
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
