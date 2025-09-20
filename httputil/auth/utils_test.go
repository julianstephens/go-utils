package auth_test

import (
	"testing"

	"github.com/julianstephens/go-utils/httputil/auth"
)

func TestHashPassword(t *testing.T) {
	password := "my_secure_password"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	
	if hash == password {
		t.Errorf("Hash should not be the same as the password")
	}
	if len(hash) == 0 {
		t.Errorf("Hash should not be empty")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "my_secure_password"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	
	if !auth.CheckPasswordHash(password, hash) {
		t.Errorf("CheckPasswordHash should return true for correct password")
	}
	if auth.CheckPasswordHash("wrong_password", hash) {
		t.Errorf("CheckPasswordHash should return false for incorrect password")
	}
}