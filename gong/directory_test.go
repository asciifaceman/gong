package gong

import "testing"

func TestEncode(t *testing.T) {
	d := &DirectoryEntryHeader{
		ID:          []byte("exampleID"),
		FNAME:       []byte("exampleFile"),
		FileType:    []byte("txt"),
		Offset:      1234,
		Size:        5678,
		Compression: 0,
	}

	encoded := d.Encode()
	if len(encoded) == 0 {
		t.Error("Encoding failed, got empty slice")
		return
	}

	expectedSize := d.EncodedSize()
	if uint32(len(encoded)) != expectedSize {
		t.Errorf("Encoded size mismatch: got %d, want %d", len(encoded), expectedSize)
	}

}
