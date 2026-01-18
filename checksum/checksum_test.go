package checksum

import (
	"testing"
)

// TestCRC32C verifies CRC32-C computation
func TestCRC32C(t *testing.T) {
	tests := []struct {
		data []byte
		name string
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "single byte",
			data: []byte{0x00},
		},
		{
			name: "simple string",
			data: []byte("hello"),
		},
		{
			name: "longer data",
			data: []byte("The quick brown fox jumps over the lazy dog"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify computation is consistent and returns uint32
			result1 := CRC32C(tt.data)
			result2 := CRC32C(tt.data)

			// Should be deterministic
			if result1 != result2 {
				t.Errorf("CRC32C should be deterministic")
			}

			// For non-empty data, should usually be non-zero (though technically possible to be 0)
			_ = result1 // Just verify it computes without error
		})
	}
}

// TestCRC32IEEE verifies standard CRC32 computation
func TestCRC32IEEE(t *testing.T) {
	data := []byte("hello")
	crc := CRC32IEEE(data)

	// IEEE CRC32 should produce a valid uint32
	if crc == 0 && len(data) > 0 {
		t.Errorf("CRC32IEEE(%q) should not be 0 for non-empty data", data)
	}
}

// TestVerifyCRC32C verifies the verification function
func TestVerifyCRC32C(t *testing.T) {
	data := []byte("test data")
	checksum := CRC32C(data)

	if !VerifyCRC32C(data, checksum) {
		t.Errorf("VerifyCRC32C should return true for valid checksum")
	}

	if VerifyCRC32C(data, checksum+1) {
		t.Errorf("VerifyCRC32C should return false for invalid checksum")
	}
}

// TestAppendCRC32C verifies checksum appending
func TestAppendCRC32C(t *testing.T) {
	data := []byte("test data")
	withChecksum := AppendCRC32C(data)

	// Should have 4 more bytes
	if len(withChecksum) != len(data)+4 {
		t.Errorf("AppendCRC32C should append 4 bytes, got %d", len(withChecksum)-len(data))
	}

	// Should be verifiable
	if !VerifyCRC32CWithSelf(withChecksum) {
		t.Errorf("VerifyCRC32CWithSelf should verify appended checksum")
	}
}

// TestStripAndVerifyCRC32C verifies stripping and verification
func TestStripAndVerifyCRC32C(t *testing.T) {
	originalData := []byte("test data")
	withChecksum := AppendCRC32C(originalData)

	strippedData, verified := StripAndVerifyCRC32C(withChecksum)

	if !verified {
		t.Errorf("StripAndVerifyCRC32C should verify correct checksum")
	}

	if string(strippedData) != string(originalData) {
		t.Errorf("StripAndVerifyCRC32C should return original data")
	}
}

// TestStripAndVerifyInvalidChecksum tests with corrupted data
func TestStripAndVerifyInvalidChecksum(t *testing.T) {
	data := AppendCRC32C([]byte("test data"))

	// Corrupt the data
	data[0] ^= 0xFF

	_, verified := StripAndVerifyCRC32C(data)
	if verified {
		t.Errorf("StripAndVerifyCRC32C should fail for corrupted data")
	}
}

// TestVerifyCRC32CWithSelfInvalid tests self-verification with bad checksum
func TestVerifyCRC32CWithSelfInvalid(t *testing.T) {
	data := []byte("test data")
	withChecksum := AppendCRC32C(data)

	// Corrupt the checksum bytes
	withChecksum[len(withChecksum)-1] ^= 0xFF

	if VerifyCRC32CWithSelf(withChecksum) {
		t.Errorf("VerifyCRC32CWithSelf should fail for corrupted checksum")
	}
}

// TestVerifyCRC32CWithSelfShortData tests with data too short
func TestVerifyCRC32CWithSelfShortData(t *testing.T) {
	data := []byte("abc") // Only 3 bytes, need at least 4 for checksum

	if VerifyCRC32CWithSelf(data) {
		t.Errorf("VerifyCRC32CWithSelf should fail for data shorter than 4 bytes")
	}
}

// TestNewCRC32CWriter tests incremental hashing
func TestNewCRC32CWriter(t *testing.T) {
	data := []byte("hello world")

	// Compute directly
	directCRC := CRC32C(data)

	// Compute incrementally
	h := NewCRC32CWriter()
	h.Write([]byte("hello"))
	h.Write([]byte(" "))
	h.Write([]byte("world"))
	incrementalCRC := h.Sum32()

	if directCRC != incrementalCRC {
		t.Errorf("Incremental CRC32C should match direct computation: 0x%x != 0x%x", incrementalCRC, directCRC)
	}
}

// TestNewCRC32IEEEWriter tests IEEE CRC32 incremental hashing
func TestNewCRC32IEEEWriter(t *testing.T) {
	data := []byte("test")

	directCRC := CRC32IEEE(data)

	h := NewCRC32IEEEWriter()
	h.Write(data)
	incrementalCRC := h.Sum32()

	if directCRC != incrementalCRC {
		t.Errorf("Incremental IEEE CRC32 should match direct computation")
	}
}

// TestConsistency verifies that multiple calls produce same result
func TestConsistency(t *testing.T) {
	data := []byte("consistency test data")

	crc1 := CRC32C(data)
	crc2 := CRC32C(data)
	crc3 := CRC32C(data)

	if crc1 != crc2 || crc2 != crc3 {
		t.Errorf("CRC32C should be consistent: 0x%x, 0x%x, 0x%x", crc1, crc2, crc3)
	}
}

// TestDifferentData verifies different inputs produce different checksums
func TestDifferentData(t *testing.T) {
	data1 := []byte("data1")
	data2 := []byte("data2")

	crc1 := CRC32C(data1)
	crc2 := CRC32C(data2)

	if crc1 == crc2 {
		t.Errorf("Different data should produce different checksums")
	}
}

// TestEmptyData verifies empty data handling
func TestEmptyData(t *testing.T) {
	empty := []byte{}

	crc := CRC32C(empty)
	if crc != 0 {
		t.Errorf("CRC32C of empty data should be 0, got 0x%x", crc)
	}

	verified := VerifyCRC32C(empty, 0)
	if !verified {
		t.Errorf("VerifyCRC32C of empty data with checksum 0 should be true")
	}
}
