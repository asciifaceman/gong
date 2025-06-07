package assets

type Asset struct {
	ID          string
	Filename    string
	Offset      uint32 // Offset in the gong file
	Size        uint32 // Size of the asset data
	Compression byte   // Compression type (0 = none for now)
}
