package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"crypto/sha256"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

var (
	// ErrInvalidKeySize is returned when the key size is invalid
	ErrInvalidKeySize = errors.New("invalid key size")
	// ErrInvalidCiphertext is returned when the ciphertext is invalid
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	// ErrDecryptionFailed is returned when decryption fails
	ErrDecryptionFailed = errors.New("decryption failed")
)

// AES-GCM Encryption/Decryption

// Encrypt encrypts plaintext using AES-GCM with the provided key.
// Returns the encrypted data with nonce prepended.
func Encrypt(key []byte, plaintext []byte) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, ErrInvalidKeySize
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-GCM with the provided key.
// Expects the nonce to be prepended to the ciphertext.
func Decrypt(key []byte, ciphertext []byte) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, ErrInvalidKeySize
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidCiphertext
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

// PBKDF2 Key Derivation

// DeriveKey derives a key from a password using PBKDF2 with SHA-256.
// saltSize specifies the size of the salt in bytes (recommended: 16 or 32).
// iterations specifies the number of iterations (recommended: 100000 or higher).
// keyLen specifies the desired key length in bytes.
func DeriveKey(password string, saltSize, iterations, keyLen int) (key []byte, salt []byte, err error) {
	salt = make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	key = pbkdf2.Key([]byte(password), salt, iterations, keyLen, sha256.New)
	return key, salt, nil
}

// DeriveKeyWithSalt derives a key from a password and existing salt using PBKDF2 with SHA-256.
func DeriveKeyWithSalt(password string, salt []byte, iterations, keyLen int) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, keyLen, sha256.New)
}

// Random Key Generation

// GenerateRandomKey generates a cryptographically secure random key of the specified length.
func GenerateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	return key, nil
}

// GenerateAESKey generates a random AES key of the specified size (16, 24, or 32 bytes).
func GenerateAESKey(keySize int) ([]byte, error) {
	if keySize != 16 && keySize != 24 && keySize != 32 {
		return nil, ErrInvalidKeySize
	}
	return GenerateRandomKey(keySize)
}

// Bcrypt Password Hashing

// HashPassword hashes a plaintext password using bcrypt with the default cost.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// HashPasswordWithCost hashes a plaintext password using bcrypt with the specified cost.
// Cost should be between bcrypt.MinCost and bcrypt.MaxCost.
func HashPasswordWithCost(password string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return "", fmt.Errorf("invalid cost: must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// VerifyPassword compares a plaintext password with a bcrypt hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Constant-time Secure Comparison

// SecureCompare performs a constant-time comparison of two byte slices.
// Returns true if they are equal, false otherwise.
func SecureCompare(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// SecureCompareString performs a constant-time comparison of two strings.
// Returns true if they are equal, false otherwise.
func SecureCompareString(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// Base64 Encoding/Decoding

// EncodeBase64 encodes data to base64 string using standard encoding.
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes a base64 string using standard encoding.
func DecodeBase64(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return decoded, nil
}

// EncodeBase64URL encodes data to base64 string using URL-safe encoding.
func EncodeBase64URL(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// DecodeBase64URL decodes a base64 string using URL-safe encoding.
func DecodeBase64URL(encoded string) ([]byte, error) {
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 URL: %w", err)
	}
	return decoded, nil
}
