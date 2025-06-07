package gong

import "errors"

var (
	ErrFileExists   = errors.New("file already exists")
	ErrFileNotFound = errors.New("file not found")
	ErrFileEmpty    = errors.New("file is empty")
	ErrFileRead     = errors.New("failed to read file")
	ErrFileOpen     = errors.New("failed to open file")
	ErrFileCreate   = errors.New("failed to create file")
	ErrFileStat     = errors.New("failed to stat file")
	ErrFileIsDir    = errors.New("file is a directory, not a regular file")

	ErrInvalidGongFile = errors.New("invalid gong file format")
)
