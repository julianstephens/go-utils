package checksum

// VerifyCRC32C verifies that data matches the expected CRC32-C checksum.
// Returns true if the checksum matches, false otherwise.
func VerifyCRC32C(data []byte, expectedChecksum uint32) bool {
	return CRC32C(data) == expectedChecksum
}

// VerifyCRC32IEEE verifies that data matches the expected IEEE CRC32 checksum.
func VerifyCRC32IEEE(data []byte, expectedChecksum uint32) bool {
	return CRC32IEEE(data) == expectedChecksum
}

// VerifyCRC32Koopman verifies that data matches the expected Koopman CRC32 checksum.
func VerifyCRC32Koopman(data []byte, expectedChecksum uint32) bool {
	return CRC32Koopman(data) == expectedChecksum
}

// VerifyCRC32CWithSelf verifies that data contains its own CRC32-C checksum at the end.
// The last 4 bytes of data are treated as the expected checksum (little-endian).
// Returns true if the checksum matches.
func VerifyCRC32CWithSelf(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	// Extract last 4 bytes as checksum (little-endian)
	payload := data[:len(data)-4]
	expectedChecksum := uint32(data[len(data)-4]) |
		uint32(data[len(data)-3])<<8 |
		uint32(data[len(data)-2])<<16 |
		uint32(data[len(data)-1])<<24

	return VerifyCRC32C(payload, expectedChecksum)
}

// AppendCRC32C computes the CRC32-C checksum of data and appends it (little-endian) to data.
// Returns the data with checksum appended.
func AppendCRC32C(data []byte) []byte {
	crc := CRC32C(data)

	// Append checksum in little-endian format
	return append(data,
		byte(crc),
		byte(crc>>8),
		byte(crc>>16),
		byte(crc>>24),
	)
}

// StripAndVerifyCRC32C removes the last 4 bytes (assumed to be a CRC32-C checksum)
// and verifies it matches the data. Returns the data without checksum and verification result.
func StripAndVerifyCRC32C(data []byte) ([]byte, bool) {
	if len(data) < 4 {
		return nil, false
	}

	payload := data[:len(data)-4]
	verified := VerifyCRC32CWithSelf(data)

	return payload, verified
}
