# Checksum Package

The `checksum` package provides fast cryptographic checksum utilities for data integrity verification. It includes CRC32 variants optimized for different use cases.

## Features

- **CRC32-C (Castagnoli)**: Fast checksums for WAL records and storage
- **CRC32 (IEEE)**: Standard IEEE 802.3 checksums
- **CRC32 (Koopman)**: Alternative polynomial variant
- **Verification**: Data integrity checks with multiple approaches
- **Streaming**: Incremental hash computation for large data
- **Self-checksumming**: Inline checksum append/verify/strip

## Installation

```bash
go get github.com/julianstephens/go-utils/checksum
```

## Quick Start

```go
import "github.com/julianstephens/go-utils/checksum"

// Compute CRC32-C checksum
data := []byte("hello world")
crc := checksum.CRC32C(data)

// Verify checksum
if checksum.VerifyCRC32C(data, crc) {
	fmt.Println("Data is intact")
}

// Append checksum to data
dataWithCRC := checksum.AppendCRC32C(data)

// Verify and strip checksum
original, verified := checksum.StripAndVerifyCRC32C(dataWithCRC)
```

## API Reference

**Direct Checksums:**
- `CRC32C(data []byte) uint32` - CRC32-C (Castagnoli) - recommended for WAL records
- `CRC32IEEE(data []byte) uint32` - Standard IEEE 802.3 CRC32
- `CRC32Koopman(data []byte) uint32` - Koopman polynomial variant

**Verification:**
- `VerifyCRC32C(data []byte, expected uint32) bool` - Verify CRC32-C checksum
- `VerifyCRC32IEEE(data []byte, expected uint32) bool` - Verify IEEE CRC32 checksum
- `VerifyCRC32Koopman(data []byte, expected uint32) bool` - Verify Koopman checksum

**Self-Checksumming:**
- `AppendCRC32C(data []byte) []byte` - Append CRC32-C checksum to data (little-endian)
- `VerifyCRC32CWithSelf(data []byte) bool` - Verify checksum in last 4 bytes
- `StripAndVerifyCRC32C(data []byte) ([]byte, bool)` - Remove and verify checksum

**Streaming:**
- `NewCRC32CWriter() hash.Hash32` - Incremental CRC32-C hasher
- `NewCRC32IEEEWriter() hash.Hash32` - Incremental IEEE CRC32 hasher
- `NewCRC32KoopmanWriter() hash.Hash32` - Incremental Koopman hasher

## Use Cases

**WAL Record Integrity:**
```go
// Append checksum to WAL record and later verify
record := []byte{/* record data */}
recordWithChecksum := checksum.AppendCRC32C(record)

originalRecord, ok := checksum.StripAndVerifyCRC32C(recordWithChecksum)
if !ok {
	fmt.Println("Record corrupted")
}
```

**Large File Streaming:**
```go
h := checksum.NewCRC32CWriter()
for {
	chunk := readChunk()
	if chunk == nil {break}
	h.Write(chunk)
}
crc := h.Sum32()
```

**Pre-computed Checksums:**
```go
data := []byte("important data")
crc := checksum.CRC32C(data)
store(data, crc)

// Later, verify
storedData, storedCRC := retrieve()
if !checksum.VerifyCRC32C(storedData, storedCRC) {
	fmt.Println("Data corrupted")
}
```

## CRC32 Variants

| Variant | Polynomial | Use Case | Speed |
|---------|-----------|----------|-------|
| **CRC32-C** | 0x1EDC6F41 | WAL records, storage systems | Very Fast |
| **CRC32 (IEEE)** | 0x04C11DB7 | General purpose, ZIP/Ethernet | Fast |
| **Koopman** | 0x741B8CD7 | Network protocols, error detection | Fast |

## Performance Notes

- CRC32 operations are extremely fast (< 1µs for typical records)
- Pre-computed lookup tables cached for each variant
- Streaming computation uses same algorithm as direct computation
- Streaming is more memory-efficient for large data

## Testing

Run tests with:

```bash
go test ./checksum
```

Run tests with race detection:

```bash
go test -race ./checksum
```

## Limitations

- **Non-cryptographic**: Not suitable for security (use SHA256 instead)
- **Bit-flip patterns**: Sensitive to certain corruption patterns
- **Little-endian**: Appended checksums use little-endian byte order
- **Untrusted parties**: Not suitable for authentication

## Related Packages

- `github.com/julianstephens/go-utils/filelock` - File locking for exclusive access
- `github.com/julianstephens/go-utils/helpers` - General file system utilities
