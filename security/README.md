# Security Package

The `security` package provides cryptographic utilities including AES-GCM encryption/decryption, PBKDF2 key derivation, HKDF key derivation, random key generation, bcrypt password hashing and verification, constant-time secure comparison, and base64 encoding/decoding.

## Features

- **AES-GCM Encryption/Decryption**: Authenticated encryption with AES-128, AES-192, and AES-256
- **PBKDF2 Key Derivation**: Secure key derivation from passwords using PBKDF2 with SHA-256
- **HKDF Key Derivation**: HMAC-based key derivation function for generating cryptographically independent keys
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

Authenticated encryption with AES-128, AES-192, and AES-256.

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

Secure key derivation from passwords using PBKDF2 with SHA-256.

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

### HKDF Key Derivation

Derive multiple independent keys from a master key using HKDF.

```go
package main

import (
    "log"
    
    "github.com/julianstephens/go-utils/security"
)

func main() {
    masterKey, _ := security.GenerateRandomKey(32)
    
    // Derive two independent keys for different purposes
    accessKey, refreshKey, _ := security.DeriveKeyPair(
        masterKey, 
        "access-salt", "refresh-salt",
        "JWT-access", "JWT-refresh",
        32,
    )
    
    _ = accessKey
    _ = refreshKey
}
```

### Random Key Generation

Generate cryptographically secure random keys.

```go
package main

import (
    "github.com/julianstephens/go-utils/security"
)

func main() {
    key128, _ := security.GenerateAESKey(16)   // AES-128
    key256, _ := security.GenerateAESKey(32)   // AES-256
    randomKey, _ := security.GenerateRandomKey(64)
    
    _ = key128
    _ = key256
    _ = randomKey
}
```

### Bcrypt Password Hashing

Secure password hashing and verification.

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

Compare sensitive data safely, resistant to timing attacks.

```go
package main

import (
    "github.com/julianstephens/go-utils/security"
)

func main() {
    secret1 := []byte("api_key")
    secret2 := []byte("api_key")
    
    if security.SecureCompare(secret1, secret2) {
        // Keys match
    }
    
    if security.SecureCompareString("token", "token") {
        // Tokens match
    }
}
```

### Base64 Encoding/Decoding

Encode and decode with standard or URL-safe base64.

```go
package main

import (
    "github.com/julianstephens/go-utils/security"
)

func main() {
    data := []byte("Hello, World!")
    
    encoded := security.EncodeBase64(data)
    decoded, _ := security.DecodeBase64(encoded)
    
    urlEncoded := security.EncodeBase64URL(data)
    urlDecoded, _ := security.DecodeBase64URL(urlEncoded)
    
    _ = decoded
    _ = urlDecoded
}
```

### Complete Example: Secure Data Storage

Combine multiple features for secure encrypted storage:

```go
package main

import (
    "log"
    
    "github.com/julianstephens/go-utils/security"
)

type SecureStorage struct {
    salt []byte
    pass string
}

func New(password string) (*SecureStorage, error) {
    salt, _ := security.GenerateRandomKey(16)
    return &SecureStorage{salt: salt, pass: password}, nil
}

func (s *SecureStorage) Encrypt(data []byte) (string, error) {
    key := security.DeriveKeyWithSalt(s.pass, s.salt, 100000, 32)
    ciphertext, err := security.Encrypt(key, data)
    if err != nil {
        return "", err
    }
    return security.EncodeBase64(ciphertext), nil
}

func (s *SecureStorage) Decrypt(encoded string) ([]byte, error) {
    ciphertext, _ := security.DecodeBase64(encoded)
    key := security.DeriveKeyWithSalt(s.pass, s.salt, 100000, 32)
    return security.Decrypt(key, ciphertext)
}

func main() {
    storage, _ := New("password123")
    encrypted, _ := storage.Encrypt([]byte("secret data"))
    decrypted, _ := storage.Decrypt(encrypted)
    _ = decrypted
}
```

## API Reference

### AES-GCM Functions

- `Encrypt(key []byte, plaintext []byte) ([]byte, error)` — Encrypt data using AES-GCM
- `Decrypt(key []byte, ciphertext []byte) ([]byte, error)` — Decrypt data using AES-GCM

### Key Derivation Functions

**PBKDF2:**
- `DeriveKey(password string, saltSize, iterations, keyLen int) (key []byte, salt []byte, err error)` — Derive key with new random salt
- `DeriveKeyWithSalt(password string, salt []byte, iterations, keyLen int) []byte` — Derive key with existing salt

**HKDF:**
- `DeriveKeyPair(masterKey []byte, salt1, salt2, info1, info2 string, keyLength int) (key1, key2 []byte, err error)` — Derive two independent keys
- `DeriveKeyHKDF(masterKey []byte, salt, info string, keyLength int) ([]byte, error)` — Derive single key using HKDF

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
- Store keys securely (environment variables, key management systems)
- Use strong master passwords for derivation
- Rotate keys regularly

### AES-GCM
- Uses unique random nonces per operation
- Provides confidentiality and authenticity
- Good for data at rest and in transit

### PBKDF2
- Minimum 100,000 iterations (password-based)
- Use 16+ byte salt (32 bytes recommended)

### HKDF
- Ideal for deriving multiple independent keys
- Use different salt/info for different purposes ("JWT-access" vs "JWT-refresh")

### Bcrypt
- Default cost (10) suitable for most applications
- Higher costs (12-15) for better security but slower

### Timing Attacks
- **Always** use `SecureCompare` for sensitive comparisons
- Never use `==` or `bytes.Equal` for secrets/tokens/hashes

## Best Practices

1. Use AES-256 (32-byte keys) for high-security applications
2. Validate input data lengths and formats
3. Never ignore cryptographic errors
4. Use strong, long master passwords for key derivation
5. Store salts alongside encrypted data (salts don't need to be secret)
6. Use HKDF for deriving multiple keys (don't reuse same key)
7. Use `SecureCompare` for all sensitive comparisons
8. Review cryptographic code regularly

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

// Used by JWT authentication in httputil/auth package
import "github.com/julianstephens/go-utils/httputil/auth"

// The JWT package uses security.DeriveKeyPair for generating 
// separate access and refresh token signing keys
jwtManager, err := auth.NewJWTManager("secret", time.Hour, "issuer")
```