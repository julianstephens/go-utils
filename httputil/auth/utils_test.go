package auth_test

import (
	"testing"

	"github.com/julianstephens/go-utils/httputil/auth"
	tst "github.com/julianstephens/go-utils/tests"
)

func TestHashPassword(t *testing.T) {
	password := "my_secure_password"
	hash, err := auth.HashPassword(password)
	tst.AssertNoError(t, err)

	tst.AssertFalse(t, hash == password, "Hash should not be the same as the password")
	tst.AssertTrue(t, len(hash) != 0, "Hash should not be empty")
}

func TestCheckPasswordHash(t *testing.T) {
	password := "my_secure_password"
	hash, err := auth.HashPassword(password)
	tst.AssertNoError(t, err)

	tst.AssertTrue(t, auth.CheckPasswordHash(password, hash), "CheckPasswordHash should return true for correct password")
	tst.AssertFalse(t, auth.CheckPasswordHash("wrong_password", hash), "CheckPasswordHash should return false for incorrect password")
}
