package gong

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/asciifaceman/hobocode"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/afero"
)

func InitializeGongFile(f afero.File) error {
	hobocode.Iinfof(2, "Initializing new gong file %s...", f.Name())

	d := DirectoryHeader{
		DirCount: 0,
		DirSize:  0,
	}

	if err := SOG.Write(f); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write SOG preamble: %w", err))
	}
	if err := GONG.Write(f); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write GONG preamble: %w", err))
	}
	if err := d.Write(f); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write directory header: %w", err))
	}
	if _, err := f.Write([]byte{NULL}); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write NULL byte after directory header: %w", err))
	}
	if err := EOG.Write(f); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write EOG preamble: %w", err))
	}

	return nil
}

func LoadGongfile(f afero.Fs, path string) (*GongFile, error) {
	hobocode.Infof("Loading gong file %s...", path)

	file, err := AcquireFile(f, path)
	if err != nil {
		return nil, err
	}

	if !SOG.IsAt(file, 0) {
		return nil, errors.Join(ErrInvalidGongFile, errors.New("file does not start with SOG preamble"))
	}
	if !GONG.IsAt(file, int64(len(SOG))) {
		return nil, errors.Join(ErrInvalidGongFile, errors.New("file does not contain GONG preamble after SOG"))
	}

	d, err := ReadDirectoryHeader(file)
	if err != nil {
		return nil, errors.Join(ErrInvalidGongFile, err)
	}

	gf := &GongFile{
		Path:    path,
		File:    file,
		D:       d,
		Packed:  make([]*Asset, 0, d.DirCount),
		Pending: make([]*Asset, 0),
	}

	if err := gf.LoadDirectory(); err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to load directory from gong file %s", path))
	}

	return gf, nil
}

type GongFile struct {
	Path            string          // Path to the gong file
	File            afero.File      // Open file handle
	D               DirectoryHeader // Directory header containing directory count and size
	Packed          []*Asset        // Packed assets in the gong file
	PackedEndOffset int64           // Offset where the packed assets end
	Pending         []*Asset        // Assets pending to be written to the gong file
}

func (g *GongFile) AssetPacked(asset *Asset) bool {
	for _, a := range g.Packed {
		if a.ID == asset.ID {
			return true
		}
	}
	return false
}

func (g *GongFile) DirectorySize() uint32 {
	size := uint32(0)
	for _, asset := range g.Packed {
		size += asset.EncodedSize()
	}
	return size
}

func (g *GongFile) LoadDirectory() error {
	hobocode.Info("Loading directory from gong file...")

	g.Packed = make([]*Asset, 0, g.D.DirCount)
	offset := AssetStartOffset

	if _, err := g.File.Seek(offset, io.SeekStart); err != nil {
		return errors.Join(err, fmt.Errorf("failed to seek to directory entries in gong file %s", g.Path))
	}

	maxOffset := offset // Track the maximum offset for packed assets

	for i := range g.D.DirCount {
		idLen := make([]byte, 1)
		if _, err := g.File.Read(idLen); err != nil {
			return errors.Join(err, fmt.Errorf("failed to read ID length for entry %d in gong file %s", i, g.Path))
		}
		idBytes := make([]byte, idLen[0])
		if _, err := g.File.Read(idBytes); err != nil {
			return errors.Join(err, fmt.Errorf("failed to read ID for entry %d in gong file %s", i, g.Path))
		}

		nameLen := make([]byte, 1)
		if _, err := g.File.Read(nameLen); err != nil {
			return errors.Join(err, fmt.Errorf("failed to read name length for entry %d in gong file %s", i, g.Path))
		}
		nameBytes := make([]byte, nameLen[0])
		if _, err := g.File.Read(nameBytes); err != nil {
			return errors.Join(err, fmt.Errorf("failed to read name for entry %d in gong file %s", i, g.Path))
		}

		typeLen := make([]byte, 1)
		if _, err := g.File.Read(typeLen); err != nil {
			return errors.Join(err, fmt.Errorf("failed to read type length for entry %d in gong file %s", i, g.Path))
		}
		typeBytes := make([]byte, typeLen[0])
		if _, err := g.File.Read(typeBytes); err != nil {
			return errors.Join(err, fmt.Errorf("failed to read type for entry %d in gong file %s", i, g.Path))
		}

		// size & compression

		var (
			aoffset     uint32
			size        uint32
			compression byte
		)

		if err := binary.Read(g.File, binary.BigEndian, &offset); err != nil {
			return errors.Join(err, fmt.Errorf("failed to read offset for entry %d in gong file %s", i, g.Path))
		}
		if err := binary.Read(g.File, binary.BigEndian, &size); err != nil {
			return errors.Join(err, fmt.Errorf("failed to read size for entry %d in gong file %s", i, g.Path))
		}
		if err := binary.Read(g.File, binary.BigEndian, &compression); err != nil {
			return errors.Join(err, fmt.Errorf("failed to read compression for entry %d in gong file %s", i, g.Path))
		}

		g.Packed = append(g.Packed, &Asset{
			ID:          string(idBytes),
			Filename:    string(nameBytes),
			FileType:    string(typeBytes),
			Offset:      aoffset,
			Size:        size,
			Compression: compression,
		})

		end := offset + int64(size)
		if end > maxOffset {
			maxOffset = end
		}
	}

	g.PackedEndOffset = maxOffset

	spew.Dump(g.Packed)

	for i, asset := range g.Packed {
		// read the asset data
		if _, err := g.File.Seek(int64(offset), io.SeekStart); err != nil {
			return errors.Join(err, fmt.Errorf("failed to seek to asset data for entry %d in gong file %s", i, g.Path))
		}
		data := make([]byte, asset.Size)
		if _, err := g.File.Read(data); err != nil {
			spew.Dump(asset)
			return errors.Join(err, fmt.Errorf("failed to read asset data for entry %d in gong file %s", i, g.Path))
		}
		asset.Content = data
	}

	return nil
}

func (g *GongFile) AppendAssets(f afero.Fs, assets map[string]string) error {
	if len(assets) == 0 {
		return fmt.Errorf("no assets to append to gong file %s", g.Path)
	}

	for id, source := range assets {
		asset, err := NewAssetFromFile(f, id, source)
		if err != nil {
			return fmt.Errorf("failed to create asset from file %s: %w", source, err)
		}

		if err := g.Append(asset); err != nil {
			return fmt.Errorf("failed to append asset %s to gong file %s: %w", asset.ID, g.Path, err)
		}
		hobocode.Iinfof(2, "Asset %s staged for packing in gong file %s", asset.ID, g.Path)
	}

	return nil
}

func (g *GongFile) Append(asset *Asset) error {
	if g.AssetPacked(asset) {
		return fmt.Errorf("asset %s is already packed in gong file %s", asset.ID, g.Path)
	}

	if asset.Size == 0 {
		return fmt.Errorf("asset %s has size 0, cannot append to gong file %s", asset.ID, g.Path)
	}
	if asset.Offset != 0 {
		return fmt.Errorf("asset %s already has an offset set, cannot append to gong file %s", asset.ID, g.Path)
	}

	g.Pending = append(g.Pending, asset)
	return nil
}

// GetAsset returns the []byte data of an asset from the gong file by its ID.
func (g *GongFile) GetAsset(id string) (*Asset, []byte, error) {
	for _, asset := range g.Packed {
		if asset.ID == id {
			assetOffset := int64(asset.Offset)
			if _, err := g.File.Seek(assetOffset, io.SeekStart); err != nil {
				return nil, nil, errors.Join(err, fmt.Errorf("failed to seek to asset %s in gong file %s", id, g.Path))
			}
			data := make([]byte, asset.Size)
			if _, err := g.File.Read(data); err != nil {
				return nil, nil, errors.Join(err, fmt.Errorf("failed to read asset %s from gong file %s", id, g.Path))
			}
			if len(data) == 0 {
				return nil, nil, fmt.Errorf("asset %s in gong file %s is empty", id, g.Path)
			}

			return asset, data, nil
		}
	}
	return nil, nil, fmt.Errorf("asset %s not found in gong file %s", id, g.Path)
}

// GetShadowFile returns a file handle to an in-memory copy of the gong file
// that can be used for writing. It reads the entire gong file into memory.
func (g *GongFile) GetShadowFile(f afero.Fs) (afero.File, error) {
	if _, err := g.File.Seek(0, io.SeekStart); err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to seek to start of gong file %s", g.Path))
	}
	var dat []byte
	if _, err := g.File.Read(dat); err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to read gong file %s", g.Path))
	}

	m, err := AcquireFile(f, g.Path)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to acquire in-memory gong file %s for writing", g.Path))
	}
	if _, err := m.Write(dat); err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to write to in-memory gong file %s", g.Path))
	}
	return m, nil
}

func (g *GongFile) Write(f afero.Fs) error {
	hobocode.Info("Preparing to write gong file...")

	for _, asset := range g.Pending {
		b, err := asset.ReadAsset(f)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to read asset %s for gong file %s", asset.ID, g.Path))
		}
		asset.Content = b
	}

	g.Packed = append(g.Packed, g.Pending...)
	g.Pending = make([]*Asset, 0)

	d := DirectoryHeader{
		DirCount: uint16(len(g.Packed)),
		DirSize:  g.DirectorySize(),
	}

	var size uint32
	for _, asset := range g.Packed {
		size += asset.EncodedSize()
	}

	of := uint32(AssetStartOffset)
	for _, asset := range g.Packed {
		asset.Offset = of
		of += asset.Size
	}

	hobocode.Iinfo(2, "creating shadow file for changes...")
	shadowFile, err := g.GetShadowFile(f)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to get shadow file for gong file %s", g.Path))
	}

	// write the preamble
	if err := SOG.Write(shadowFile); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write SOG preamble: %w", err))
	}
	if err := GONG.Write(shadowFile); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write GONG preamble: %w", err))
	}
	if err := d.Write(shadowFile); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write directory header: %w", err))
	}
	if _, err := shadowFile.Write([]byte{NULL}); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write NULL byte after directory header: %w", err))
	}

	// write the asset directory
	hobocode.Iinfo(2, "writing asset directory...")
	for _, asset := range g.Packed {
		enc := asset.Encode()
		spew.Dump(enc)
		if enc == nil {
			return fmt.Errorf("failed to encode asset %s for gong file %s", asset.ID, g.Path)
		}
		if _, err := shadowFile.Write(enc); err != nil {
			return errors.Join(err, fmt.Errorf("failed to write asset %s to gong file %s", asset.ID, g.Path))
		}
	}

	// write the packed assets
	hobocode.Iinfo(2, "writing packed assets...")

	if err := SOA.Write(shadowFile); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write SOA preamble: %w", err))
	}

	for i, asset := range g.Packed {
		if _, err := shadowFile.Write(asset.Content); err != nil {
			return errors.Join(err, fmt.Errorf("failed to write asset %s to gong file %s", asset.ID, g.Path))
		}
		hobocode.Iinfof(2, "Wrote asset %d/%d: %s", i+1, len(g.Packed), asset.ID)
	}

	if err := EOA.Write(shadowFile); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write EOA preamble: %w", err))
	}

	// write the EOG postamble
	if err := EOG.Write(shadowFile); err != nil {
		return errors.Join(err, fmt.Errorf("failed to write EOG postamble: %w", err))
	}

	// write the shadow file to the original file
	hobocode.Iinfo(2, "writing gong file...")
	if _, err := shadowFile.Seek(0, io.SeekStart); err != nil {
		return errors.Join(err, fmt.Errorf("failed to seek to start of shadow file for gong file %s", g.Path))
	}
	if _, err := g.File.Seek(0, io.SeekStart); err != nil {
		return errors.Join(err, fmt.Errorf("failed to seek to start of gong file %s", g.Path))
	}
	if _, err := io.Copy(g.File, shadowFile); err != nil {
		return errors.Join(err, fmt.Errorf("failed to copy shadow file to gong file %s", g.Path))
	}
	defer shadowFile.Close()
	hobocode.Isuccessf(2, "Successfully wrote gong file %s with %d assets", g.Path, len(g.Packed))
	return nil
}
