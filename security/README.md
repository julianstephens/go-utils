# Security Package

The `security` package provides cryptographic utilities including AES-GCM encryption/decryption, PBKDF2 key derivation, random key generation, bcrypt password hashing and verification, constant-time secure comparison, and base64 encoding/decoding.

## Features

- **AES-GCM Encryption/Decryption**: Authenticated encryption with AES-128, AES-192, and AES-256
- **PBKDF2 Key Derivation**: Secure key derivation from passwords using PBKDF2 with SHA-256
- **Random Key Generation**: Cryptographically secure random key generation
- **Bcrypt Password Hashing**: Secure password hashing and verification using bcrypt
- **Constant-time Comparison**: Secure comparison functions resistant to timing attacks
- **Base64 Encoding/Decoding**: Both standard and URL-safe base64 encoding/decoding

## Installation

```bash
go get github.com/julianstephens/go-utils
```

## Usage

### AES-GCM Encryption/Decryption

AES-GCM provides authenticated encryption, ensuring both confidentiality and integrity of your data.

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/julianstephens/go-utils/security"
)

func main() {
    // Generate a 32-byte key for AES-256
    key, err := security.GenerateAESKey(32)
    if err != nil {
        log.Fatalf("Failed to generate key: %v", err)
    }
    
    plaintext := []byte("This is a secret message!")
    
    // Encrypt the data
    ciphertext, err := security.Encrypt(key, plaintext)
    if err != nil {
        log.Fatalf("Failed to encrypt: %v", err)
    }
    
    fmt.Printf("Encrypted data length: %d bytes\n", len(ciphertext))
    
    // Decrypt the data
    decrypted, err := security.Decrypt(key, ciphertext)
    if err != nil {
        log.Fatalf("Failed to decrypt: %v", err)
    }
    
    fmt.Printf("Decrypted: %s\n", string(decrypted))
}
```

### PBKDF2 Key Derivation

Derive cryptographic keys from passwords using PBKDF2 with secure parameters.

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/julianstephens/go-utils/security"
)

func main() {
    password := "user_password123"
    
    // Derive a key with a random salt
    key, salt, err := security.DeriveKey(password, 16, 100000, 32)
    if err != nil {
        log.Fatalf("Failed to derive key: %v", err)
    }
    
    fmt.Printf("Derived key: %x\n", key)
    fmt.Printf("Salt: %x\n", salt)
    
    // Later, derive the same key using the stored salt
    sameKey := security.DeriveKeyWithSalt(password, salt, 100000, 32)
    fmt.Printf("Same key: %x\n", sameKey)
    
    // Verify they match
    if security.SecureCompare(key, sameKey) {
        fmt.Println("Keys match!")
    }
}
```

### Random Key Generation

Generate cryptographically secure random keys for various purposes.

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/julianstephens/go-utils/security"
)

func main() {
    // Generate different AES key sizes
    aes128Key, err := security.GenerateAESKey(16) // AES-128
    if err != nil {
        log.Fatalf("Failed to generate AES-128 key: %v", err)
    }
    
    aes256Key, err := security.GenerateAESKey(32) // AES-256
    if err != nil {
        log.Fatalf("Failed to generate AES-256 key: %v", err)
    }
    
    // Generate a random key of any length
    randomKey, err := security.GenerateRandomKey(64)
    if err != nil {
        log.Fatalf("Failed to generate random key: %v", err)
    }
    
    fmt.Printf("AES-128 key: %x\n", aes128Key)
    fmt.Printf("AES-256 key: %x\n", aes256Key)
    fmt.Printf("64-byte random key: %x\n", randomKey)
}
```

### Bcrypt Password Hashing

Securely hash and verify passwords using bcrypt.

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/julianstephens/go-utils/security"
)

func main() {
    password := "user_password123"
    
    // Hash the password with default cost
    hash, err := security.HashPassword(password)
    if err != nil {
        log.Fatalf("Failed to hash password: %v", err)
    }
    
    fmt.Printf("Password hash: %s\n", hash)
    
    // Verify the password
    if security.VerifyPassword(password, hash) {
        fmt.Println("Password is correct!")
    } else {
        fmt.Println("Password is incorrect!")
    }
    
    // Hash with custom cost (higher cost = more secure but slower)
    highCostHash, err := security.HashPasswordWithCost(password, 12)
    if err != nil {
        log.Fatalf("Failed to hash password: %v", err)
    }
    
    fmt.Printf("High-cost hash: %s\n", highCostHash)
}
```

### Constant-time Secure Comparison

Perform secure comparisons that are resistant to timing attacks.

```go
package main

import (
    "fmt"
    
    "github.com/julianstephens/go-utils/security"
)

func main() {
    // Compare byte slices securely
    secret1 := []byte("secret_api_key_12345")
    secret2 := []byte("secret_api_key_12345")
    secret3 := []byte("different_secret_key")
    
    fmt.Printf("secret1 == secret2: %t\n", security.SecureCompare(secret1, secret2))
    fmt.Printf("secret1 == secret3: %t\n", security.SecureCompare(secret1, secret3))
    
    // Compare strings securely
    token1 := "abc123def456"
    token2 := "abc123def456"
    token3 := "different_token"
    
    fmt.Printf("token1 == token2: %t\n", security.SecureCompareString(token1, token2))
    fmt.Printf("token1 == token3: %t\n", security.SecureCompareString(token1, token3))
}
```

### Base64 Encoding/Decoding

Encode and decode data using both standard and URL-safe base64 encoding.

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/julianstephens/go-utils/security"
)

func main() {
    data := []byte("Hello, World! This contains special chars: +/=")
    
    // Standard base64 encoding
    stdEncoded := security.EncodeBase64(data)
    fmt.Printf("Standard base64: %s\n", stdEncoded)
    
    stdDecoded, err := security.DecodeBase64(stdEncoded)
    if err != nil {
        log.Fatalf("Failed to decode standard base64: %v", err)
    }
    fmt.Printf("Decoded: %s\n", string(stdDecoded))
    
    // URL-safe base64 encoding (good for URLs and filenames)
    urlEncoded := security.EncodeBase64URL(data)
    fmt.Printf("URL-safe base64: %s\n", urlEncoded)
    
    urlDecoded, err := security.DecodeBase64URL(urlEncoded)
    if err != nil {
        log.Fatalf("Failed to decode URL-safe base64: %v", err)
    }
    fmt.Printf("Decoded: %s\n", string(urlDecoded))
}
```

### Complete Example: Secure Data Storage

Here's a complete example showing how to combine multiple features for secure data storage:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/julianstephens/go-utils/security"
)

type SecureStorage struct {
    masterPassword string
    salt          []byte
}

func NewSecureStorage(masterPassword string) (*SecureStorage, error) {
    // Generate a random salt for key derivation
    salt, err := security.GenerateRandomKey(16)
    if err != nil {
        return nil, fmt.Errorf("failed to generate salt: %w", err)
    }
    
    return &SecureStorage{
        masterPassword: masterPassword,
        salt:          salt,
    }, nil
}

func (s *SecureStorage) deriveKey() []byte {
    return security.DeriveKeyWithSalt(s.masterPassword, s.salt, 100000, 32)
}

func (s *SecureStorage) EncryptData(data []byte) (string, error) {
    // Derive encryption key from master password
    key := s.deriveKey()
    
    // Encrypt the data
    ciphertext, err := security.Encrypt(key, data)
    if err != nil {
        return "", fmt.Errorf("encryption failed: %w", err)
    }
    
    // Encode to base64 for storage
    encoded := security.EncodeBase64(ciphertext)
    return encoded, nil
}

func (s *SecureStorage) DecryptData(encodedData string) ([]byte, error) {
    // Decode from base64
    ciphertext, err := security.DecodeBase64(encodedData)
    if err != nil {
        return nil, fmt.Errorf("failed to decode base64: %w", err)
    }
    
    // Derive decryption key from master password
    key := s.deriveKey()
    
    // Decrypt the data
    plaintext, err := security.Decrypt(key, ciphertext)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %w", err)
    }
    
    return plaintext, nil
}

func (s *SecureStorage) GetSaltForStorage() string {
    return security.EncodeBase64(s.salt)
}

func main() {
    // Create secure storage with master password
    storage, err := NewSecureStorage("my_master_password_123")
    if err != nil {
        log.Fatalf("Failed to create secure storage: %v", err)
    }
    
    // Encrypt sensitive data
    sensitiveData := []byte("Credit card: 1234-5678-9012-3456, SSN: 123-45-6789")
    encrypted, err := storage.EncryptData(sensitiveData)
    if err != nil {
        log.Fatalf("Failed to encrypt data: %v", err)
    }
    
    fmt.Printf("Encrypted data: %s\n", encrypted)
    fmt.Printf("Salt (store this): %s\n", storage.GetSaltForStorage())
    
    // Decrypt the data
    decrypted, err := storage.DecryptData(encrypted)
    if err != nil {
        log.Fatalf("Failed to decrypt data: %v", err)
    }
    
    fmt.Printf("Decrypted data: %s\n", string(decrypted))
}
```

## API Reference

### AES-GCM Functions

- `Encrypt(key []byte, plaintext []byte) ([]byte, error)` — Encrypt data using AES-GCM
- `Decrypt(key []byte, ciphertext []byte) ([]byte, error)` — Decrypt data using AES-GCM

### Key Derivation Functions

- `DeriveKey(password string, saltSize, iterations, keyLen int) (key []byte, salt []byte, err error)` — Derive key with new random salt
- `DeriveKeyWithSalt(password string, salt []byte, iterations, keyLen int) []byte` — Derive key with existing salt

### Random Key Generation

- `GenerateRandomKey(length int) ([]byte, error)` — Generate random key of specified length
- `GenerateAESKey(keySize int) ([]byte, error)` — Generate AES key (16, 24, or 32 bytes)

### Password Hashing Functions

- `HashPassword(password string) (string, error)` — Hash password with default cost
- `HashPasswordWithCost(password string, cost int) (string, error)` — Hash password with custom cost
- `VerifyPassword(password, hash string) bool` — Verify password against hash

### Secure Comparison Functions

- `SecureCompare(a, b []byte) bool` — Constant-time comparison of byte slices
- `SecureCompareString(a, b string) bool` — Constant-time comparison of strings

### Base64 Functions

- `EncodeBase64(data []byte) string` — Standard base64 encoding
- `DecodeBase64(encoded string) ([]byte, error)` — Standard base64 decoding
- `EncodeBase64URL(data []byte) string` — URL-safe base64 encoding
- `DecodeBase64URL(encoded string) ([]byte, error)` — URL-safe base64 decoding

## Error Types

The package defines several error constants:

- `ErrInvalidKeySize` — Invalid AES key size (must be 16, 24, or 32 bytes)
- `ErrInvalidCiphertext` — Invalid ciphertext format
- `ErrDecryptionFailed` — Decryption failed (wrong key or corrupted data)

## Security Considerations

### Key Management

- Store encryption keys securely (consider using environment variables or key management systems)
- Use strong master passwords for key derivation
- Regularly rotate encryption keys when possible

### AES-GCM

- Each encryption operation uses a unique random nonce
- Provides both confidentiality and authenticity
- Suitable for encrypting data at rest and in transit

### PBKDF2

- Use at least 100,000 iterations for password-based key derivation
- Use a minimum 16-byte salt (32 bytes recommended)
- Consider using Argon2 for new applications (not included in this package)

### Bcrypt

- Default cost (10) provides good security for most applications
- Higher costs (12-15) provide better security but slower performance
- Cost should be adjusted based on hardware capabilities and security requirements

### Timing Attacks

- Always use `SecureCompare` functions when comparing sensitive data
- Never use `==` or `bytes.Equal` for comparing secrets, tokens, or hashes

## Best Practices

1. **Use appropriate key sizes**: AES-256 (32-byte keys) for high-security applications
2. **Validate input data**: Check data lengths and formats before processing
3. **Handle errors properly**: Don't ignore cryptographic operation errors
4. **Use strong passwords**: For key derivation, use long, complex passwords
5. **Store salts safely**: Keep salts alongside encrypted data, they don't need to be secret
6. **Regular security audits**: Review cryptographic implementations regularly

## Thread Safety

All functions in this package are thread-safe and can be called concurrently from multiple goroutines.

## Integration

Works well with other go-utils packages:

```go
// Use with logger for security event logging
import "github.com/julianstephens/go-utils/logger"

log := logger.New(logger.Config{Level: "info"})
hash, err := security.HashPassword(password)
if err != nil {
    log.WithError(err).Error("Failed to hash password")
}

// Use with helpers for configuration
import "github.com/julianstephens/go-utils/helpers"

keySize := helpers.Default(configKeySize, 32)
key, err := security.GenerateAESKey(keySize)
```