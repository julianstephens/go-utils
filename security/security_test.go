package security_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/julianstephens/go-utils/security"
	tst "github.com/julianstephens/go-utils/tests"
	"golang.org/x/crypto/bcrypt"
)

// Test AES-GCM Encryption/Decryption

func TestEncryptDecrypt(t *testing.T) {
	// Test with 32-byte key (AES-256)
	key, err := security.GenerateAESKey(32)
	tst.RequireNoError(t, err)

	plaintext := []byte("Hello, World! This is a secret message.")

	// Encrypt
	ciphertext, err := security.Encrypt(key, plaintext)
	tst.RequireNoError(t, err)
	tst.AssertTrue(t, len(ciphertext) > len(plaintext), "Ciphertext should be longer than plaintext")

	// Decrypt
	decrypted, err := security.Decrypt(key, ciphertext)
	tst.RequireNoError(t, err)
	tst.AssertDeepEqual(t, decrypted, plaintext)
}

func TestEncryptDecryptDifferentKeySizes(t *testing.T) {
	keySizes := []int{16, 24, 32} // AES-128, AES-192, AES-256
	plaintext := []byte("Test message for different key sizes")

	for _, keySize := range keySizes {
		t.Run(fmt.Sprintf("KeySize%d", keySize*8), func(t *testing.T) {
			key, err := security.GenerateAESKey(keySize)
			tst.RequireNoError(t, err)

			ciphertext, err := security.Encrypt(key, plaintext)
			tst.RequireNoError(t, err)

			decrypted, err := security.Decrypt(key, ciphertext)
			tst.RequireNoError(t, err)
			tst.AssertDeepEqual(t, decrypted, plaintext)
		})
	}
}

func TestEncryptWithInvalidKeySize(t *testing.T) {
	invalidKey := make([]byte, 15) // Invalid key size
	plaintext := []byte("test")

	_, err := security.Encrypt(invalidKey, plaintext)
	tst.AssertTrue(t, err == security.ErrInvalidKeySize, "Should return ErrInvalidKeySize")
}

func TestDecryptWithInvalidKeySize(t *testing.T) {
	invalidKey := make([]byte, 15) // Invalid key size
	ciphertext := []byte("test")

	_, err := security.Decrypt(invalidKey, ciphertext)
	tst.AssertTrue(t, err == security.ErrInvalidKeySize, "Should return ErrInvalidKeySize")
}

func TestDecryptWithInvalidCiphertext(t *testing.T) {
	key, err := security.GenerateAESKey(32)
	tst.RequireNoError(t, err)

	// Ciphertext too short (less than nonce size)
	shortCiphertext := []byte("short")
	_, err = security.Decrypt(key, shortCiphertext)
	tst.AssertTrue(t, err == security.ErrInvalidCiphertext, "Should return ErrInvalidCiphertext")
}

func TestDecryptWithWrongKey(t *testing.T) {
	key1, err := security.GenerateAESKey(32)
	tst.RequireNoError(t, err)
	key2, err := security.GenerateAESKey(32)
	tst.RequireNoError(t, err)

	plaintext := []byte("secret message")

	// Encrypt with key1
	ciphertext, err := security.Encrypt(key1, plaintext)
	tst.RequireNoError(t, err)

	// Try to decrypt with key2
	_, err = security.Decrypt(key2, ciphertext)
	tst.AssertTrue(t, err == security.ErrDecryptionFailed, "Should return ErrDecryptionFailed")
}

func TestEncryptDecryptEmptyData(t *testing.T) {
	key, err := security.GenerateAESKey(32)
	tst.RequireNoError(t, err)

	plaintext := []byte("")

	ciphertext, err := security.Encrypt(key, plaintext)
	tst.RequireNoError(t, err)

	decrypted, err := security.Decrypt(key, ciphertext)
	tst.RequireNoError(t, err)

	// Handle nil vs empty slice comparison
	if len(plaintext) == 0 && len(decrypted) == 0 {
		// Both are effectively empty, this is correct
		return
	}
	tst.AssertDeepEqual(t, decrypted, plaintext)
}

// Test PBKDF2 Key Derivation

func TestDeriveKey(t *testing.T) {
	password := "my_secure_password"
	saltSize := 16
	iterations := 100000
	keyLen := 32

	key, salt, err := security.DeriveKey(password, saltSize, iterations, keyLen)
	tst.RequireNoError(t, err)

	tst.AssertTrue(t, len(key) == keyLen, "Key length should match requested length")
	tst.AssertTrue(t, len(salt) == saltSize, "Salt length should match requested length")
}

func TestDeriveKeyWithSalt(t *testing.T) {
	password := "my_secure_password"
	salt := []byte("fixed_salt_16byt")
	iterations := 100000
	keyLen := 32

	key1 := security.DeriveKeyWithSalt(password, salt, iterations, keyLen)
	key2 := security.DeriveKeyWithSalt(password, salt, iterations, keyLen)

	tst.AssertTrue(t, len(key1) == keyLen, "Key length should match requested length")
	tst.AssertDeepEqual(t, key1, key2)
}

func TestDeriveKeyDifferentPasswords(t *testing.T) {
	salt := []byte("fixed_salt_16byt")
	iterations := 100000
	keyLen := 32

	key1 := security.DeriveKeyWithSalt("password1", salt, iterations, keyLen)
	key2 := security.DeriveKeyWithSalt("password2", salt, iterations, keyLen)

	tst.AssertFalse(t, bytes.Equal(key1, key2), "Different passwords should produce different keys")
}

func TestDeriveKeyDifferentSalts(t *testing.T) {
	password := "same_password"
	iterations := 100000
	keyLen := 32

	key1, salt1, err := security.DeriveKey(password, 16, iterations, keyLen)
	tst.RequireNoError(t, err)
	key2, salt2, err := security.DeriveKey(password, 16, iterations, keyLen)
	tst.RequireNoError(t, err)

	tst.AssertFalse(t, bytes.Equal(salt1, salt2), "Different calls should produce different salts")
	tst.AssertFalse(t, bytes.Equal(key1, key2), "Different salts should produce different keys")
}

// Test Random Key Generation

func TestGenerateRandomKey(t *testing.T) {
	lengths := []int{16, 24, 32, 64}

	for _, length := range lengths {
		t.Run(fmt.Sprintf("Length%d", length), func(t *testing.T) {
			key, err := security.GenerateRandomKey(length)
			tst.RequireNoError(t, err)
			tst.AssertTrue(t, len(key) == length, "Key length should match requested length")

			// Generate another key to ensure they're different
			key2, err := security.GenerateRandomKey(length)
			tst.RequireNoError(t, err)
			tst.AssertFalse(t, bytes.Equal(key, key2), "Two random keys should be different")
		})
	}
}

func TestGenerateAESKey(t *testing.T) {
	validSizes := []int{16, 24, 32}

	for _, size := range validSizes {
		t.Run(fmt.Sprintf("Size%d", size), func(t *testing.T) {
			key, err := security.GenerateAESKey(size)
			tst.RequireNoError(t, err)
			tst.AssertTrue(t, len(key) == size, "Key length should match requested size")
		})
	}
}

func TestGenerateAESKeyInvalidSize(t *testing.T) {
	invalidSizes := []int{8, 15, 17, 31, 33, 64}

	for _, size := range invalidSizes {
		t.Run(fmt.Sprintf("InvalidSize%d", size), func(t *testing.T) {
			_, err := security.GenerateAESKey(size)
			tst.AssertTrue(t, err == security.ErrInvalidKeySize, "Should return ErrInvalidKeySize for invalid size")
		})
	}
}

// Test Bcrypt Password Hashing

func TestHashPassword(t *testing.T) {
	password := "my_secure_password"
	hash, err := security.HashPassword(password)
	tst.RequireNoError(t, err)

	tst.AssertFalse(t, hash == password, "Hash should not be the same as the password")
	tst.AssertTrue(t, len(hash) > 0, "Hash should not be empty")
	tst.AssertTrue(t, strings.HasPrefix(hash, "$2a$"), "Hash should have bcrypt prefix")
}

func TestHashPasswordWithCost(t *testing.T) {
	password := "my_secure_password"
	cost := 12

	hash, err := security.HashPasswordWithCost(password, cost)
	tst.RequireNoError(t, err)

	tst.AssertFalse(t, hash == password, "Hash should not be the same as the password")
	tst.AssertTrue(t, len(hash) > 0, "Hash should not be empty")

	// Verify the cost is correct by checking if the hash can be verified
	tst.AssertTrue(t, security.VerifyPassword(password, hash), "Password should verify against its hash")
}

func TestHashPasswordWithInvalidCost(t *testing.T) {
	password := "test"
	invalidCosts := []int{bcrypt.MinCost - 1, bcrypt.MaxCost + 1}

	for _, cost := range invalidCosts {
		t.Run(fmt.Sprintf("Cost%d", cost), func(t *testing.T) {
			_, err := security.HashPasswordWithCost(password, cost)
			tst.AssertTrue(t, err != nil, "Should return error for invalid cost")
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "my_secure_password"
	hash, err := security.HashPassword(password)
	tst.RequireNoError(t, err)

	tst.AssertTrue(t, security.VerifyPassword(password, hash), "VerifyPassword should return true for correct password")
	tst.AssertFalse(t, security.VerifyPassword("wrong_password", hash), "VerifyPassword should return false for incorrect password")
}

func TestVerifyPasswordDifferentHashes(t *testing.T) {
	password := "same_password"

	hash1, err := security.HashPassword(password)
	tst.RequireNoError(t, err)
	hash2, err := security.HashPassword(password)
	tst.RequireNoError(t, err)

	// Hashes should be different due to random salt
	tst.AssertFalse(t, hash1 == hash2, "Different hash calls should produce different hashes")

	// But both should verify the same password
	tst.AssertTrue(t, security.VerifyPassword(password, hash1), "Password should verify against first hash")
	tst.AssertTrue(t, security.VerifyPassword(password, hash2), "Password should verify against second hash")
}

// Test Constant-time Secure Comparison

func TestSecureCompare(t *testing.T) {
	a := []byte("secret_value")
	b := []byte("secret_value")
	c := []byte("different_value")

	tst.AssertTrue(t, security.SecureCompare(a, b), "Identical byte slices should be equal")
	tst.AssertFalse(t, security.SecureCompare(a, c), "Different byte slices should not be equal")
}

func TestSecureCompareDifferentLengths(t *testing.T) {
	a := []byte("short")
	b := []byte("much_longer_string")

	tst.AssertFalse(t, security.SecureCompare(a, b), "Different length byte slices should not be equal")
}

func TestSecureCompareEmpty(t *testing.T) {
	a := []byte("")
	b := []byte("")
	c := []byte("not_empty")

	tst.AssertTrue(t, security.SecureCompare(a, b), "Empty byte slices should be equal")
	tst.AssertFalse(t, security.SecureCompare(a, c), "Empty and non-empty byte slices should not be equal")
}

func TestSecureCompareString(t *testing.T) {
	a := "secret_value"
	b := "secret_value"
	c := "different_value"

	tst.AssertTrue(t, security.SecureCompareString(a, b), "Identical strings should be equal")
	tst.AssertFalse(t, security.SecureCompareString(a, c), "Different strings should not be equal")
}

func TestSecureCompareStringDifferentLengths(t *testing.T) {
	a := "short"
	b := "much_longer_string"

	tst.AssertFalse(t, security.SecureCompareString(a, b), "Different length strings should not be equal")
}

func TestSecureCompareStringEmpty(t *testing.T) {
	a := ""
	b := ""
	c := "not_empty"

	tst.AssertTrue(t, security.SecureCompareString(a, b), "Empty strings should be equal")
	tst.AssertFalse(t, security.SecureCompareString(a, c), "Empty and non-empty strings should not be equal")
}

// Test Base64 Encoding/Decoding

func TestEncodeDecodeBase64(t *testing.T) {
	data := []byte("Hello, World! This is a test message for base64 encoding.")

	encoded := security.EncodeBase64(data)
	tst.AssertTrue(t, len(encoded) > 0, "Encoded string should not be empty")

	decoded, err := security.DecodeBase64(encoded)
	tst.RequireNoError(t, err)
	tst.AssertDeepEqual(t, decoded, data)
}

func TestEncodeDecodeBase64URL(t *testing.T) {
	data := []byte("Hello, World! This is a test message for base64 URL encoding.")

	encoded := security.EncodeBase64URL(data)
	tst.AssertTrue(t, len(encoded) > 0, "Encoded string should not be empty")

	decoded, err := security.DecodeBase64URL(encoded)
	tst.RequireNoError(t, err)
	tst.AssertDeepEqual(t, decoded, data)
}

func TestEncodeDecodeBase64Empty(t *testing.T) {
	data := []byte("")

	encoded := security.EncodeBase64(data)
	decoded, err := security.DecodeBase64(encoded)
	tst.RequireNoError(t, err)
	tst.AssertDeepEqual(t, decoded, data)
}

func TestDecodeBase64InvalidInput(t *testing.T) {
	invalidInputs := []string{
		"invalid base64!",
		"SGVsbG8gV29ybGQ==invalid",
		"not base64 at all",
	}

	for _, input := range invalidInputs {
		t.Run("Invalid_"+input[:min(10, len(input))], func(t *testing.T) {
			_, err := security.DecodeBase64(input)
			tst.AssertTrue(t, err != nil, "Should return error for invalid base64")
		})
	}
}

func TestDecodeBase64URLInvalidInput(t *testing.T) {
	invalidInputs := []string{
		"invalid base64!",
		"SGVsbG8gV29ybGQ==invalid",
		"not base64 at all",
	}

	for _, input := range invalidInputs {
		t.Run("Invalid_"+input[:min(10, len(input))], func(t *testing.T) {
			_, err := security.DecodeBase64URL(input)
			tst.AssertTrue(t, err != nil, "Should return error for invalid base64 URL")
		})
	}
}

func TestBase64VsBase64URL(t *testing.T) {
	// Data that will contain URL-unsafe characters when base64 encoded
	data := []byte("subjects? and objects>")

	stdEncoded := security.EncodeBase64(data)
	urlEncoded := security.EncodeBase64URL(data)

	// They should be different due to URL-safe encoding
	tst.AssertFalse(t, stdEncoded == urlEncoded, "Standard and URL encodings should be different for data with special chars")

	// Both should decode to the same original data
	stdDecoded, err := security.DecodeBase64(stdEncoded)
	tst.RequireNoError(t, err)
	urlDecoded, err := security.DecodeBase64URL(urlEncoded)
	tst.RequireNoError(t, err)

	tst.AssertDeepEqual(t, stdDecoded, data)
	tst.AssertDeepEqual(t, urlDecoded, data)
	tst.AssertDeepEqual(t, stdDecoded, urlDecoded)
}

// Integration Tests

func TestEncryptDecryptWithDerivedKey(t *testing.T) {
	password := "user_password"
	plaintext := []byte("This is a secret message encrypted with a derived key.")

	// Derive key
	key, salt, err := security.DeriveKey(password, 16, 100000, 32)
	tst.RequireNoError(t, err)

	// Encrypt
	ciphertext, err := security.Encrypt(key, plaintext)
	tst.RequireNoError(t, err)

	// Derive the same key again
	key2 := security.DeriveKeyWithSalt(password, salt, 100000, 32)
	tst.AssertDeepEqual(t, key, key2)

	// Decrypt
	decrypted, err := security.Decrypt(key2, ciphertext)
	tst.RequireNoError(t, err)
	tst.AssertDeepEqual(t, decrypted, plaintext)
}

func TestBase64EncodeRandomData(t *testing.T) {
	randomData, err := security.GenerateRandomKey(64)
	tst.RequireNoError(t, err)

	encoded := security.EncodeBase64(randomData)
	decoded, err := security.DecodeBase64(encoded)
	tst.RequireNoError(t, err)
	tst.AssertDeepEqual(t, decoded, randomData)
}

// Test HKDF Key Derivation

func TestDeriveKeyPair(t *testing.T) {
	masterKey := []byte("super-secret-master-key-for-testing")
	salt1 := "salt1"
	salt2 := "salt2"
	info1 := "key1 info"
	info2 := "key2 info"
	keyLength := 32

	key1, key2, err := security.DeriveKeyPair(masterKey, salt1, salt2, info1, info2, keyLength)
	tst.RequireNoError(t, err)

	tst.AssertTrue(t, len(key1) == keyLength, "Key1 should have correct length")
	tst.AssertTrue(t, len(key2) == keyLength, "Key2 should have correct length")
	tst.AssertFalse(t, bytes.Equal(key1, key2), "Derived keys should be different")

	// Test deterministic behavior - same inputs should produce same outputs
	key1b, key2b, err := security.DeriveKeyPair(masterKey, salt1, salt2, info1, info2, keyLength)
	tst.RequireNoError(t, err)
	tst.AssertDeepEqual(t, key1, key1b)
	tst.AssertDeepEqual(t, key2, key2b)
}

func TestDeriveKeyPairDifferentSalts(t *testing.T) {
	masterKey := []byte("master-key")
	keyLength := 32

	// Same salts should produce same keys
	key1a, key2a, err := security.DeriveKeyPair(masterKey, "salt1", "salt2", "info1", "info2", keyLength)
	tst.RequireNoError(t, err)
	key1b, key2b, err := security.DeriveKeyPair(masterKey, "salt1", "salt2", "info1", "info2", keyLength) 
	tst.RequireNoError(t, err)
	tst.AssertDeepEqual(t, key1a, key1b)
	tst.AssertDeepEqual(t, key2a, key2b)

	// Different salts should produce different keys
	key1c, key2c, err := security.DeriveKeyPair(masterKey, "different1", "different2", "info1", "info2", keyLength)
	tst.RequireNoError(t, err)
	tst.AssertFalse(t, bytes.Equal(key1a, key1c), "Different salts should produce different key1")
	tst.AssertFalse(t, bytes.Equal(key2a, key2c), "Different salts should produce different key2")
}

func TestDeriveKeyHKDF(t *testing.T) {
	masterKey := []byte("master-key-for-single-derivation")
	salt := "unique-salt"
	info := "key info"
	keyLength := 32

	key, err := security.DeriveKeyHKDF(masterKey, salt, info, keyLength)
	tst.RequireNoError(t, err)
	tst.AssertTrue(t, len(key) == keyLength, "Key should have correct length")

	// Test deterministic behavior
	key2, err := security.DeriveKeyHKDF(masterKey, salt, info, keyLength)
	tst.RequireNoError(t, err)
	tst.AssertDeepEqual(t, key, key2)

	// Different salt should produce different key
	key3, err := security.DeriveKeyHKDF(masterKey, "different-salt", info, keyLength)
	tst.RequireNoError(t, err)
	tst.AssertFalse(t, bytes.Equal(key, key3), "Different salt should produce different key")
}

func TestDeriveKeyHKDFDifferentLengths(t *testing.T) {
	masterKey := []byte("master-key")
	salt := "salt"
	info := "info"

	lengths := []int{16, 24, 32, 64}
	for _, length := range lengths {
		t.Run(fmt.Sprintf("Length%d", length), func(t *testing.T) {
			key, err := security.DeriveKeyHKDF(masterKey, salt, info, length)
			tst.RequireNoError(t, err)
			tst.AssertTrue(t, len(key) == length, "Key should have requested length")
		})
	}
}

// Helper function for min (since Go 1.21+ has this built-in, but maintaining compatibility)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
