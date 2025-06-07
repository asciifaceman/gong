package gong

import (
	"fmt"

	"github.com/asciifaceman/gong/gong/magic"
	"github.com/spf13/afero"
)

type Asset struct {
	ID          string
	Filename    string
	Offset      uint32 // Offset in the gong file
	Size        uint32 // Size of the asset data
	Compression byte   // Compression type (0 = none for now)
}

func BuildAssetEntries(files map[string]string) ([]*Asset, error) {
	assets := make([]*Asset, 0, len(files))

	f := afero.NewOsFs()

	for id, file := range files {
		info, err := f.Stat(file)
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
// and updates the Offset field of each Asset. It assumes that the assets
// slice is already populated with Asset objects and that the offsets
// are to be calculated based on the total size of the directory entries.
// It returns the total size of the directory entries.
func AssignAssetOffsets(assets []*Asset) uint32 {
	size := CalculateDirSize(assets)

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
