package assets

import (
	"fmt"
	"os"

	"github.com/asciifaceman/gong/gong/magic"
)

func BuildAssetEntries(files map[string]string) ([]*Asset, error) {
	assets := make([]*Asset, 0, len(files))

	for id, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file %s: %w", file, err)
		}
		if info.IsDir() {
			return nil, fmt.Errorf("file %s is a directory, not a regular file", file)
		}

		assets = append(assets, &Asset{
			ID:          id,
			Filename:    file,
			Offset:      0, // Offset will be set later when writing to the gong file
			Size:        uint32(info.Size()),
			Compression: 0, // Compression type (0 = none for now)
		})
	}

	return assets, nil
}

// AssignAssetOffsets calculates the offsets for each asset in the gong file
// and returns the total size of the asset entries for the directory.
func AssignAssetOffsets(assets []*Asset) uint32 {
	var size uint32

	for _, asset := range assets {
		if asset == nil {
			continue
		}
		size += 1 + uint32(len(asset.ID)) +
			1 + uint32(len(asset.Filename)) +
			4 + 4 + 1 // offset, size, and compression type
	}

	start := magic.HeaderSize + size
	current := start
	for _, asset := range assets {
		if asset == nil {
			continue
		}
		asset.Offset = current
		current += asset.Size
	}

	return size
}
