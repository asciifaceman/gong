package gong

func uint16ToBytes(value uint16) []byte {
	return []byte{byte(value >> 8), byte(value)}
}
func uint32ToBytes(value uint32) []byte {
	return []byte{byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)}
}
func bytesToUint16(b []byte) uint16 {
	if len(b) < 2 {
		return 0 // Not enough bytes to convert
	}
	return uint16(b[0])<<8 | uint16(b[1])
}
func bytesToUint32(b []byte) uint32 {
	if len(b) < 4 {
		return 0 // Not enough bytes to convert
	}
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}
