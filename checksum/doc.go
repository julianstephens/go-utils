// Package checksum provides cryptographic checksum utilities for data integrity.
//
// It includes:
//
// - CRC32-C (Castagnoli): Fast polynomial-based checksum for storage systems
// - CRC32 (standard): IEEE 802.3 polynomial for general use
// - Verification helpers: Quick data integrity checks
//
// Example usage:
//
//	// Compute CRC32-C checksum
//	crc := checksum.CRC32C(data)
//
//	// Verify data against checksum
//	ok := checksum.VerifyCRC32C(data, crc)
package checksum
