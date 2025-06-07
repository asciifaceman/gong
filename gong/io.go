package gong

import (
	"errors"
	"fmt"
	"os"

	"github.com/asciifaceman/hobocode"
	"github.com/spf13/afero"
)

// AcquireFile attempts to acquire a file at the specified path.
// If the file exists, it opens it for reading and writing.
// If the file does not exist, it creates a new file at that path.
func AcquireFile(f afero.Fs, path string) (afero.File, error) {
	if exists, err := afero.Exists(f, path); err != nil {
		return nil, errors.Join(ErrFileStat, err)
	} else if exists {
		file, err := f.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			return nil, errors.Join(ErrFileOpen, err)
		}
		return file, nil
	} else {
		hobocode.Iinfof(2, "gongfile %s does not exist, creating new file and initializing...", path)
		file, err := f.Create(path)
		if err != nil {
			return nil, errors.Join(ErrFileCreate, err)
		}
		if err := InitializeGongFile(file); err != nil {
			file.Close()
			return nil, errors.Join(err, fmt.Errorf("failed to initialize gong file %s", path))
		}
		return file, nil
	}

}

// ReadBytes reads the contents of a file and returns it as a byte slice.
func ReadBytes(f afero.Fs, path string) ([]byte, error) {
	file, err := f.Open(path)
	if err != nil {
		return nil, errors.Join(ErrFileOpen, err)
	}
	defer file.Close()
	data, err := afero.ReadAll(file)
	if err != nil {
		return nil, errors.Join(ErrFileRead, err)
	}
	if len(data) == 0 {
		return nil, ErrFileEmpty
	}
	return data, nil
}
