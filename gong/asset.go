package gong

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

// NewAssetFromFile creates a new Asset from a file in the filesystem.
// It reads the file's metadata and initializes the Asset fields.
func NewAssetFromFile(f afero.Fs, id, source string) (*Asset, error) {
	if id == "" {
		return nil, errors.New("asset ID cannot be empty")
	}
	if len(id) > 255 {
		return nil, fmt.Errorf("asset ID '%s' exceeds maximum length of 255 characters", id)
	}
	if source == "" {
		return nil, errors.New("asset source cannot be empty")
	}

	a := &Asset{
		ID:       id,
		Source:   source,
		Filename: filepath.Base(source), // Extract filename from source path
	}

	info, err := f.Stat(source)
	if err != nil {
		return nil, errors.Join(ErrFileStat, err)
	}
	if info.IsDir() {
		return nil, errors.Join(ErrFileIsDir, fmt.Errorf("file %s is a directory", source))
	}
	if info.Size() == 0 {
		return nil, ErrFileEmpty
	}
	a.Size = uint32(info.Size())

	suffix := filepath.Ext(a.Filename)
	a.FileType = suffix

	return a, nil
}

// Asset is a file bundled in a gong file.
//
// Asset header format:
//
// SOA (2 bytes) | ID length (1 byte) | ID bytes | Filename length (1 byte) |
// Filename bytes | FileType length (1 byte) | FileType bytes |
// Offset (4 bytes) | Size (4 bytes) | Compression type (1 byte) |
// EOA (2 bytes)
type Asset struct {
	ID          string // Unique identifier for the asset
	Filename    string // Name of the file sans path in the gong archive
	FileType    string // filetype of the asset by extension, e.g. "png", "ogg", etc.
	Source      string // Source path of the asset
	Content     []byte // Content of the asset, read from Source
	Offset      uint32 // Offset in the gong file
	Size        uint32 // Size of the asset data
	Compression byte   // Compression type (0 = none for now)
}

// EncodedSize returns the size of the asset in bytes accounting for metadata.
func (a *Asset) EncodedSize() uint32 {
	return 1 + // ID length byte
		uint32(len(a.ID)) + // ID bytes
		1 + // Filename length byte
		uint32(len(a.Filename)) + // Filename bytes
		1 + // FileType length byte
		uint32(len(a.FileType)) + // FileType bytes
		4 + // Offset (uint32)
		4 + // Size (uint32)
		1 // Compression type (byte)
}

// ReadAsset reads Source and returns a []byte
func (a *Asset) ReadAsset(f afero.Fs) ([]byte, error) {
	file, err := f.Open(a.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to open asset file %s: %w", a.Source, err)
	}
	defer file.Close()

	data, err := afero.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset file %s: %w", a.Source, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("asset file %s is empty", a.Source)
	}

	return data, nil
}

// SetContent
func (a *Asset) SetContent(b []byte) {
	a.Content = b
}

// Encode returns a []byte of the bundled asset of EncodedSize
func (a *Asset) Encode() []byte {
	var enc bytes.Buffer

	if _, err := enc.Write([]byte{byte(len(a.ID))}); err != nil {
		return nil // Error writing ID length, return empty slice
	}
	if _, err := enc.Write([]byte(a.ID)); err != nil {
		return nil // Error writing ID bytes, return empty slice
	}
	if _, err := enc.Write([]byte{byte(len(a.Filename))}); err != nil {
		return nil // Error writing Filename length, return empty slice
	}
	if _, err := enc.Write([]byte(a.Filename)); err != nil {
		return nil // Error writing Filename bytes, return empty slice
	}
	if _, err := enc.Write([]byte{byte(len(a.FileType))}); err != nil {
		return nil // Error writing FileType length, return empty slice
	}
	if _, err := enc.Write([]byte(a.FileType)); err != nil {
		return nil // Error writing FileType bytes, return empty slice
	}
	if _, err := enc.Write([]byte{byte(a.Offset)}); err != nil {
		return nil // Error writing Offset, return empty slice
	}
	if _, err := enc.Write([]byte{byte(a.Size)}); err != nil {
		return nil // Error writing Size, return empty slice
	}
	if _, err := enc.Write([]byte{a.Compression}); err != nil {
		return nil // Error writing Compression type, return empty slice
	}

	return enc.Bytes()
}

// BuildAssets creates a slice of Asset objects from a map of file IDs to file paths.
// It reads each file's metadata and initializes the Asset fields accordingly.
func BuildAssets(f afero.Fs, files map[string]string) ([]*Asset, error) {
	assets := make([]*Asset, 0, len(files))

	for id, source := range files {
		asset, err := NewAssetFromFile(f, id, source)
		if err != nil {
			return nil, fmt.Errorf("failed to create asset from file %s: %w", source, err)
		}
		assets = append(assets, asset)
	}

	return assets, nil
}
