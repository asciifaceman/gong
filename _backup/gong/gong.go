package gong

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/asciifaceman/gong/gong/gio"
	"github.com/asciifaceman/gong/gong/magic"
	"github.com/spf13/afero"
)

func WriteGongFile(path string, newAssets []*Asset) error {
	f := afero.NewOsFs()

	var (
		file     afero.File
		existing = make([]*Asset, 0) // existing assets in the Gong file
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

	allEntries := append(existing, newAssets...)

	// assign offsets and write asset data
	for _, entry := range allEntries {
		entry.Offset = assetDB

		dat, err := gio.ReadFileBytes(f, entry.Filename)
		if err != nil {
			return fmt.Errorf("failed to read asset file %s for entry %s: %w", entry.Filename, entry.ID, err)
		}

		if _, err := file.Seek(int64(entry.Offset), io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek to offset %d for entry %s: %w", entry.Offset, entry.ID, err)
		}

		buf := bytes.NewBuffer(dat)
		reader := bufio.NewReader(buf)
		n, err := io.CopyN(file, reader, int64(entry.Size))
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to write asset data for entry %s: %w", entry.ID, err)
		}
		if uint32(n) != entry.Size {
			return fmt.Errorf("written size %d does not match expected size %d for entry %s", n, entry.Size, entry.ID)
		}
		assetDB += entry.Size
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
