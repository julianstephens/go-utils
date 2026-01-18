package checksum

import (
	"hash"
	"hash/crc32"
)

// CRC32C computes the CRC32-C (Castagnoli) checksum of data.
// CRC32-C uses the polynomial 0x1EDC6F41 and is commonly used for
// WAL records, storage systems, and network protocols.
func CRC32C(data []byte) uint32 {
	return crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))
}

// CRC32IEEE computes the standard IEEE 802.3 CRC32 checksum.
// This is the most common CRC32 variant.
func CRC32IEEE(data []byte) uint32 {
	return crc32.Checksum(data, crc32.MakeTable(crc32.IEEE))
}

// CRC32Koopman computes the Koopman polynomial CRC32 checksum.
func CRC32Koopman(data []byte) uint32 {
	return crc32.Checksum(data, crc32.MakeTable(crc32.Koopman))
}

// NewCRC32CWriter returns a new hash.Hash32 for computing CRC32-C incrementally.
// This is useful for streaming data.
func NewCRC32CWriter() hash.Hash32 {
	return crc32.New(crc32.MakeTable(crc32.Castagnoli))
}

// NewCRC32IEEEWriter returns a new hash.Hash32 for computing IEEE CRC32 incrementally.
func NewCRC32IEEEWriter() hash.Hash32 {
	return crc32.New(crc32.MakeTable(crc32.IEEE))
}

// NewCRC32KoopmanWriter returns a new hash.Hash32 for computing Koopman CRC32 incrementally.
func NewCRC32KoopmanWriter() hash.Hash32 {
	return crc32.New(crc32.MakeTable(crc32.Koopman))
}
