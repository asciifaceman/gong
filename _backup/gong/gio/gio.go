package gio

import (
	"fmt"
	"os"

	"github.com/asciifaceman/hobocode"
	"github.com/spf13/afero"
)

var (
	// ErrFileExists is returned when the file already exists
	ErrFileExists = fmt.Errorf("file already exists")
)

// Check if file exists using the given filesystem interface
func Exists(f afero.Fs, path string) bool {
	// Check if the file exists at the given path
	if exists, err := afero.Exists(f, path); err == nil {
		return exists
	} else {
		hobocode.Errorf("Error checking if file exists: %v", err)
	}
	return false
}

// ReadFileBytes reads the contents of a file and returns it as a byte slice.
func ReadFileBytes(f afero.Fs, path string) ([]byte, error) {
	file, err := f.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()
	data, err := afero.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("file %s is empty", path)
	}
	return data, nil
}

func AcquireFile(f afero.Fs, path string) (file afero.File, err error) {
	if Exists(f, path) {
		file, err = f.OpenFile(path, os.O_RDWR, 0644)
		return
	}
	file, err = f.Create(path)
	return
}
