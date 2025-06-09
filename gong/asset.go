package gong

// Asset is directory entry for a file bundled in a gong file
//
type Asset struct {
	ID          string // Unique id for the asset
	Filename    string // Name of the file sans path and ext
	FileType    string // type of the file is the extension
	Offset      uint32 // Offset in the gong file where the content starts
	Size        uint32 // Size of the content stored at offset
	Compression byte   // Compression type (0 = none)
}
