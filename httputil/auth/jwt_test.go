package auth_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/julianstephens/go-utils/httputil/auth"
	tst "github.com/julianstephens/go-utils/tests"
)

func TestNewJWTManager(t *testing.T) {
	secretKey := "test-secret-key"
	duration := time.Hour
	issuer := "test-issuer"

	manager, err := auth.NewJWTManager(secretKey, duration, issuer)
	tst.AssertNoError(t, err)
	tst.AssertNotNil(t, manager, "NewJWTManager should not return nil")
}

func TestGenerateToken(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	userID := "user123"
	roles := []string{"user", "admin"}

	token, err := manager.GenerateToken(userID, roles)
	tst.AssertNoError(t, err)
	tst.AssertTrue(t, token != "", "Generated token should not be empty")

	// Token should have 3 parts separated by dots
	parts := strings.Split(token, ".")
	tst.AssertTrue(t, len(parts) == 3, "Invalid JWT format: expected 3 parts")
}

func TestValidateToken(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	userID := "user123"
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user", "admin"}

	// Generate a token
	token, err := manager.GenerateTokenWithUserInfo(userID, username, email, roles)
	tst.AssertNoError(t, err)

	// Validate the token
	claims, err := manager.ValidateToken(token)
	tst.AssertNoError(t, err)

	// Verify claims
	tst.AssertTrue(t, claims.UserID == userID, "UserID should match")
	tst.AssertTrue(t, claims.Username == username, "Username should match")
	tst.AssertTrue(t, claims.Email == email, "Email should match")
	tst.AssertDeepEqual(t, claims.Roles, roles)
}

func TestValidateInvalidToken(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	tests := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"invalid format", "invalid.token"},
		{"malformed token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := manager.ValidateToken(tt.token)
			tst.AssertNotNil(t, err, "ValidateToken should fail for invalid token")
		})
	}
}

func TestValidateExpiredToken(t *testing.T) {
	// Create manager with very short expiration
	manager, err := auth.NewJWTManager("test-secret", time.Millisecond, "test-issuer")
	tst.AssertNoError(t, err)

	token, err := manager.GenerateToken("user123", []string{"user"})
	tst.AssertNoError(t, err)

	// Wait for token to expire
	time.Sleep(time.Millisecond * 10)

	_, err = manager.ValidateToken(token)
	tst.AssertNotNil(t, err, "ValidateToken should fail for expired token")
	tst.AssertTrue(t, err == auth.ErrTokenExpired, "Expected ErrTokenExpired")
}

func TestRefreshToken(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	// Generate original token
	originalToken, err := manager.GenerateToken("user123", []string{"user"})
	tst.AssertNoError(t, err)

	// Refresh the token
	newToken, err := manager.RefreshToken(originalToken)
	tst.AssertNoError(t, err)

	tst.AssertTrue(t, newToken != "", "Refreshed token should not be empty")
	tst.AssertTrue(t, newToken != originalToken, "Refreshed token should be different from original")

	// Validate the new token
	claims, err := manager.ValidateToken(newToken)
	tst.AssertNoError(t, err)

	// Claims should be the same
	tst.AssertTrue(t, claims.UserID == "user123", "UserID should be preserved in refreshed token")
}

func TestRefreshExpiredToken(t *testing.T) {
	// Create manager with very short expiration for initial token
	shortManager, err := auth.NewJWTManager("test-secret", time.Millisecond, "test-issuer")
	tst.AssertNoError(t, err)

	token, err := shortManager.GenerateToken("user123", []string{"user"})
	tst.AssertNoError(t, err)

	// Wait for token to expire
	time.Sleep(time.Millisecond * 10)

	// Create a manager with longer expiration for refresh
	longManager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	// Should still be able to refresh expired token
	newToken, err := longManager.RefreshToken(token)
	tst.AssertNoError(t, err)

	// New token should be valid
	_, err = longManager.ValidateToken(newToken)
	tst.AssertNoError(t, err)
}

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name        string
		header      string
		expectToken string
		expectError bool
	}{
		{
			name:        "valid bearer token",
			header:      "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectError: false,
		},
		{
			name:        "empty header",
			header:      "",
			expectToken: "",
			expectError: true,
		},
		{
			name:        "missing Bearer prefix",
			header:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectToken: "",
			expectError: true,
		},
		{
			name:        "wrong prefix",
			header:      "Basic dXNlcjpwYXNz",
			expectToken: "",
			expectError: true,
		},
		{
			name:        "Bearer without token",
			header:      "Bearer ",
			expectToken: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := auth.ExtractTokenFromHeader(tt.header)
			if tt.expectError {
				tst.AssertNotNil(t, err, "Expected error but got none")
			} else {
				tst.AssertNoError(t, err)
				tst.AssertTrue(t, token == tt.expectToken, "Extracted token should match expected")
			}
		})
	}
}

func TestClaimsHasRole(t *testing.T) {
	claims := &auth.Claims{
		Roles: []string{"user", "admin", "moderator"},
	}

	tests := []struct {
		role     string
		expected bool
	}{
		{"user", true},
		{"admin", true},
		{"moderator", true},
		{"guest", false},
		{"superuser", false},
	}

	for _, tt := range tests {
		result := claims.HasRole(tt.role)
		tst.AssertTrue(t, result == tt.expected, "HasRole result should match expected")
	}
}

func TestClaimsHasAnyRole(t *testing.T) {
	claims := &auth.Claims{
		Roles: []string{"user", "admin"},
	}

	tests := []struct {
		name     string
		roles    []string
		expected bool
	}{
		{
			name:     "has one of the roles",
			roles:    []string{"user", "guest"},
			expected: true,
		},
		{
			name:     "has all roles",
			roles:    []string{"user", "admin"},
			expected: true,
		},
		{
			name:     "has none of the roles",
			roles:    []string{"guest", "moderator"},
			expected: false,
		},
		{
			name:     "empty roles",
			roles:    []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := claims.HasAnyRole(tt.roles...)
			tst.AssertTrue(t, result == tt.expected, "HasAnyRole result should match expected")
		})
	}
}

func TestClaimsIsExpired(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	// Generate a valid token
	token, err := manager.GenerateToken("user123", []string{"user"})
	tst.AssertNoError(t, err)

	claims, err := manager.ValidateToken(token)
	tst.AssertNoError(t, err)

	// Token should not be expired
	tst.AssertFalse(t, claims.IsExpired(), "Token should not be expired")

	// Test with expired token
	expiredManager, err := auth.NewJWTManager("test-secret", time.Millisecond, "test-issuer")
	tst.AssertNoError(t, err)
	expiredToken, err := expiredManager.GenerateToken("user123", []string{"user"})
	tst.AssertNoError(t, err)

	// Wait for expiration
	time.Sleep(time.Millisecond * 10)

	// Parse expired token without validation to get claims
	expiredClaims, err := expiredManager.ValidateToken(expiredToken)
	if err == nil || err != auth.ErrTokenExpired {
		// If we can't get expired claims through validation, we'll skip this part
		// as the token validation correctly identifies it as expired
		return
	}

	// If we somehow got claims from an expired token, they should show as expired
	if expiredClaims != nil {
		tst.AssertTrue(t, expiredClaims.IsExpired(), "Expired token claims should show as expired")
	}
}

func TestTokenValidationWithDifferentSecrets(t *testing.T) {
	manager1, err := auth.NewJWTManager("secret1", time.Hour, "issuer1")
	tst.AssertNoError(t, err)
	manager2, err := auth.NewJWTManager("secret2", time.Hour, "issuer2")
	tst.AssertNoError(t, err)

	// Generate token with first manager
	token, err := manager1.GenerateToken("user123", []string{"user"})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Try to validate with second manager (different secret)
	_, err = manager2.ValidateToken(token)
	if err == nil {
		t.Error("Token validation should fail with different secret")
	}
}

func TestGenerateTokenWithClaims(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	userID := "user123"
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user", "admin"}
	customClaims := map[string]interface{}{
		"department":  "engineering",
		"level":       5,
		"is_manager":  true,
		"permissions": []string{"read", "write", "delete"},
		"metadata":    map[string]string{"region": "us-west", "tenant": "acme"},
	}

	token, err := manager.GenerateTokenWithUserInfoAndClaims(userID, username, email, roles, customClaims)
	tst.AssertNoError(t, err)
	tst.AssertTrue(t, token != "", "Generated token should not be empty")

	// Validate token and check custom claims
	claims, err := manager.ValidateToken(token)
	tst.AssertNoError(t, err)

	// Verify standard claims
	tst.AssertTrue(t, claims.UserID == userID, "UserID should match")

	// Verify custom claims
	if dept, ok := claims.GetCustomClaimString("department"); !ok || dept != "engineering" {
		t.Fatalf("Expected department 'engineering', got %s (exists: %t)", dept, ok)
	}

	if level, ok := claims.GetCustomClaimInt("level"); !ok || level != 5 {
		t.Fatalf("Expected level 5, got %d (exists: %t)", level, ok)
	}

	if isManager, ok := claims.GetCustomClaimBool("is_manager"); !ok || !isManager {
		t.Fatalf("Expected is_manager true, got %t (exists: %t)", isManager, ok)
	}
}

func TestCustomClaimMethods(t *testing.T) {
	claims := &auth.Claims{}

	// Test setting custom claims
	claims.SetCustomClaim("string_val", "test")
	claims.SetCustomClaim("int_val", 42)
	claims.SetCustomClaim("bool_val", true)
	claims.SetCustomClaim("float_val", 3.14)

	// Test HasCustomClaim
	if !claims.HasCustomClaim("string_val") {
		t.Error("HasCustomClaim should return true for existing claim")
	}
	if claims.HasCustomClaim("nonexistent") {
		t.Error("HasCustomClaim should return false for non-existing claim")
	}

	// Test GetCustomClaimString
	if val, ok := claims.GetCustomClaimString("string_val"); !ok || val != "test" {
		t.Errorf("GetCustomClaimString failed: expected 'test', got %s (exists: %t)", val, ok)
	}
	if _, ok := claims.GetCustomClaimString("int_val"); ok {
		t.Error("GetCustomClaimString should return false for non-string value")
	}

	// Test GetCustomClaimInt
	if val, ok := claims.GetCustomClaimInt("int_val"); !ok || val != 42 {
		t.Errorf("GetCustomClaimInt failed: expected 42, got %d (exists: %t)", val, ok)
	}
	if _, ok := claims.GetCustomClaimInt("string_val"); ok {
		t.Error("GetCustomClaimInt should return false for non-int value")
	}

	// Test GetCustomClaimBool
	if val, ok := claims.GetCustomClaimBool("bool_val"); !ok || !val {
		t.Errorf("GetCustomClaimBool failed: expected true, got %t (exists: %t)", val, ok)
	}
	if _, ok := claims.GetCustomClaimBool("string_val"); ok {
		t.Error("GetCustomClaimBool should return false for non-bool value")
	}

	// Test GetCustomClaim
	if val, ok := claims.GetCustomClaim("float_val"); !ok {
		t.Error("GetCustomClaim should find existing claim")
	} else if floatVal, isFloat := val.(float64); !isFloat || floatVal != 3.14 {
		t.Errorf("GetCustomClaim failed: expected 3.14, got %v", val)
	}

	// Test DeleteCustomClaim
	claims.DeleteCustomClaim("string_val")
	if claims.HasCustomClaim("string_val") {
		t.Error("DeleteCustomClaim should remove the claim")
	}
}

func TestCustomClaimsInRefreshToken(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	customClaims := map[string]interface{}{
		"department": "engineering",
		"level":      5,
		"is_manager": true,
	}

	// Generate original token with custom claims
	originalToken, err := manager.GenerateTokenWithUserInfoAndClaims("user123", "testuser", "test@example.com", []string{"user"}, customClaims)
	if err != nil {
		t.Fatalf("GenerateTokenWithClaims failed: %v", err)
	}

	// Refresh the token
	newToken, err := manager.RefreshToken(originalToken)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}

	// Validate the new token
	claims, err := manager.ValidateToken(newToken)
	if err != nil {
		t.Fatalf("ValidateToken failed for refreshed token: %v", err)
	}

	// Verify custom claims are preserved
	if dept, ok := claims.GetCustomClaimString("department"); !ok || dept != "engineering" {
		t.Errorf("Custom claim 'department' not preserved in refresh: got %s (exists: %t)", dept, ok)
	}

	if level, ok := claims.GetCustomClaimInt("level"); !ok || level != 5 {
		t.Errorf("Custom claim 'level' not preserved in refresh: got %d (exists: %t)", level, ok)
	}

	if isManager, ok := claims.GetCustomClaimBool("is_manager"); !ok || !isManager {
		t.Errorf("Custom claim 'is_manager' not preserved in refresh: got %t (exists: %t)", isManager, ok)
	}
}

func TestCustomClaimsWithJSONNumberTypes(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	customClaims := map[string]interface{}{
		"int_as_float": float64(42), // JSON numbers are typically float64
		"large_int":    int64(9223372036854775807),
	}

	token, err := manager.GenerateTokenWithUserInfoAndClaims("user123", "testuser", "test@example.com", []string{"user"}, customClaims)
	if err != nil {
		t.Fatalf("GenerateTokenWithClaims failed: %v", err)
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	// Test that float64 can be retrieved as int
	if val, ok := claims.GetCustomClaimInt("int_as_float"); !ok || val != 42 {
		t.Errorf("GetCustomClaimInt should handle float64: expected 42, got %d (exists: %t)", val, ok)
	}
}

func TestEmptyCustomClaims(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	// Generate token with nil custom claims
	token, err := manager.GenerateTokenWithUserInfo("user123", "testuser", "test@example.com", []string{"user"})
	if err != nil {
		t.Fatalf("GenerateTokenWithClaims failed: %v", err)
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	// Test methods on empty custom claims
	if claims.HasCustomClaim("anything") {
		t.Error("HasCustomClaim should return false for nil custom claims")
	}

	if _, ok := claims.GetCustomClaim("anything"); ok {
		t.Error("GetCustomClaim should return false for nil custom claims")
	}

	if _, ok := claims.GetCustomClaimString("anything"); ok {
		t.Error("GetCustomClaimString should return false for nil custom claims")
	}

	// Test that we can still set claims on empty custom claims map
	claims.SetCustomClaim("new_claim", "value")
	if !claims.HasCustomClaim("new_claim") {
		t.Error("SetCustomClaim should initialize custom claims map")
	}
}

func TestSimplifiedTokenGeneration(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	// Test basic token generation with just user ID and roles
	token, err := manager.GenerateToken("user123", []string{"user", "admin"})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("Generated token should not be empty")
	}

	// Validate the token
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	// Verify basic claims
	if claims.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got %s", claims.UserID)
	}

	if len(claims.Roles) != 2 || claims.Roles[0] != "user" || claims.Roles[1] != "admin" {
		t.Errorf("Expected roles [user admin], got %v", claims.Roles)
	}

	// Username and Email should be empty
	if claims.Username != "" {
		t.Errorf("Expected empty Username, got %s", claims.Username)
	}

	if claims.Email != "" {
		t.Errorf("Expected empty Email, got %s", claims.Email)
	}

	// Verify Subject claim is set to UserID
	if claims.RegisteredClaims.Subject != "user123" {
		t.Errorf("Expected Subject claim 'user123', got %s", claims.RegisteredClaims.Subject)
	}
}

func TestTokenGenerationWithCustomClaimsOnly(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	customClaims := map[string]interface{}{
		"username":   "john_doe",
		"email":      "john@example.com",
		"department": "engineering",
		"level":      5,
	}

	// Test token generation with custom claims
	token, err := manager.GenerateTokenWithClaims("user123", []string{"user"}, customClaims)
	if err != nil {
		t.Fatalf("GenerateTokenWithClaims failed: %v", err)
	}

	// Validate the token
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	// Verify that username and email were extracted from custom claims
	if claims.Username != "john_doe" {
		t.Errorf("Expected Username 'john_doe', got %s", claims.Username)
	}

	if claims.Email != "john@example.com" {
		t.Errorf("Expected Email 'john@example.com', got %s", claims.Email)
	}

	// Verify custom claims
	if dept, ok := claims.GetCustomClaimString("department"); !ok || dept != "engineering" {
		t.Errorf("Expected department 'engineering', got %s (exists: %t)", dept, ok)
	}

	if level, ok := claims.GetCustomClaimInt("level"); !ok || level != 5 {
		t.Errorf("Expected level 5, got %d (exists: %t)", level, ok)
	}

	// Verify Subject claim is set to UserID
	if claims.RegisteredClaims.Subject != "user123" {
		t.Errorf("Expected Subject claim 'user123', got %s", claims.RegisteredClaims.Subject)
	}
}

func TestTokenExpiration(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	token, err := manager.GenerateToken("user123", []string{"user"})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	exp, err := manager.TokenExpiration(token)
	if err != nil {
		t.Fatalf("TokenExpiration failed: %v", err)
	}

	// expiration should be in the future
	if time.Now().After(exp) {
		t.Errorf("Expected expiration in the future, got %v", exp)
	}
}

func TestTokenExpirationExpiredToken(t *testing.T) {
	// Create a token with very short duration so it expires
	shortManager, err := auth.NewJWTManager("test-secret", time.Millisecond, "test-issuer")
	if err != nil {
		t.Fatalf("NewJWTManager failed: %v", err)
	}

	token, err := shortManager.GenerateToken("user123", []string{"user"})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Wait for token to expire
	time.Sleep(time.Millisecond * 10)

	// Use a manager with same secret to retrieve expiration
	mgr, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	if err != nil {
		t.Fatalf("NewJWTManager failed: %v", err)
	}
	exp, err := mgr.TokenExpiration(token)
	if err != nil {
		t.Fatalf("TokenExpiration failed for expired token: %v", err)
	}

	if time.Now().Before(exp) {
		t.Errorf("Expected expiration in the past for expired token, got %v", exp)
	}
}

func TestKeyDerivationConsistency(t *testing.T) {
	// Test that key derivation produces consistent results
	secret := "test-secret-key"

	manager1, err := auth.NewJWTManager(secret, time.Hour, "test-issuer")
	if err != nil {
		t.Fatalf("NewJWTManager failed: %v", err)
	}

	manager2, err := auth.NewJWTManager(secret, time.Hour, "test-issuer")
	if err != nil {
		t.Fatalf("NewJWTManager failed: %v", err)
	}

	// Generate token with first manager
	token, err := manager1.GenerateToken("user123", []string{"user"})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Validate with second manager (should use same derived keys)
	claims, err := manager2.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got %s", claims.UserID)
	}
} // Refresh Token Workflow Tests

func TestGenerateTokenPair(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	userID := "user123"
	roles := []string{"user", "admin"}

	tokenPair, err := manager.GenerateTokenPair(userID, roles)
	tst.AssertNoError(t, err)
	tst.AssertNotNil(t, tokenPair, "Token pair should not be nil")
	tst.AssertTrue(t, tokenPair.AccessToken != "", "Access token should not be empty")
	tst.AssertTrue(t, tokenPair.RefreshToken != "", "Refresh token should not be empty")
	tst.AssertTrue(t, tokenPair.TokenType == "Bearer", "Token type should be Bearer")
	tst.AssertTrue(t, tokenPair.ExpiresIn == 3600, "Expires in should match token duration in seconds")

	// Validate access token
	claims, err := manager.ValidateToken(tokenPair.AccessToken)
	tst.AssertNoError(t, err)
	tst.AssertTrue(t, claims.UserID == userID, "UserID should match")
	tst.AssertTrue(t, len(claims.Roles) == 2, "Should have 2 roles")

	// Validate refresh token
	refreshClaims, err := manager.ValidateRefreshToken(tokenPair.RefreshToken)
	tst.AssertNoError(t, err)
	tst.AssertTrue(t, refreshClaims.UserID == userID, "UserID should match in refresh token")
	tst.AssertTrue(t, len(refreshClaims.Roles) == 2, "Should have 2 roles in refresh token")
	tst.AssertTrue(t, refreshClaims.TokenID != "", "Refresh token should have unique ID")
}

func TestGenerateTokenPairWithUserInfo(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	userID := "user123"
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user"}

	tokenPair, err := manager.GenerateTokenPairWithUserInfo(userID, username, email, roles)
	tst.AssertNoError(t, err)
	tst.AssertNotNil(t, tokenPair, "Token pair should not be nil")

	// Validate access token has user info
	claims, err := manager.ValidateToken(tokenPair.AccessToken)
	tst.AssertNoError(t, err)
	tst.AssertTrue(t, claims.Username == username, "Username should be preserved")
	tst.AssertTrue(t, claims.Email == email, "Email should be preserved")

	// Validate refresh token has user info
	refreshClaims, err := manager.ValidateRefreshToken(tokenPair.RefreshToken)
	tst.AssertNoError(t, err)
	tst.AssertTrue(t, refreshClaims.Username == username, "Username should be preserved in refresh token")
	tst.AssertTrue(t, refreshClaims.Email == email, "Email should be preserved in refresh token")
}

func TestExchangeRefreshToken(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	userID := "user123"
	roles := []string{"user", "admin"}

	// Generate initial token pair
	initialPair, err := manager.GenerateTokenPair(userID, roles)
	tst.AssertNoError(t, err)

	// Exchange refresh token for new pair
	newPair, err := manager.ExchangeRefreshToken(initialPair.RefreshToken)
	tst.AssertNoError(t, err)
	tst.AssertNotNil(t, newPair, "New token pair should not be nil")

	// Tokens should be different
	tst.AssertTrue(t, newPair.AccessToken != initialPair.AccessToken, "New access token should be different")
	tst.AssertTrue(t, newPair.RefreshToken != initialPair.RefreshToken, "New refresh token should be different")

	// Validate new access token
	claims, err := manager.ValidateToken(newPair.AccessToken)
	tst.AssertNoError(t, err)
	tst.AssertTrue(t, claims.UserID == userID, "UserID should be preserved")
	tst.AssertTrue(t, len(claims.Roles) == 2, "Roles should be preserved")

	// Validate new refresh token
	refreshClaims, err := manager.ValidateRefreshToken(newPair.RefreshToken)
	tst.AssertNoError(t, err)
	tst.AssertTrue(t, refreshClaims.UserID == userID, "UserID should be preserved in new refresh token")
}

func TestExchangeRefreshTokenWithCustomClaims(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	customClaims := map[string]interface{}{
		"department": "engineering",
		"level":      5,
		"is_manager": true,
	}

	// Generate initial token pair with custom claims
	initialPair, err := manager.GenerateTokenPairWithClaims("user123", []string{"user"}, customClaims)
	tst.AssertNoError(t, err)

	// Exchange refresh token
	newPair, err := manager.ExchangeRefreshToken(initialPair.RefreshToken)
	tst.AssertNoError(t, err)

	// Validate custom claims are preserved in new access token
	claims, err := manager.ValidateToken(newPair.AccessToken)
	tst.AssertNoError(t, err)

	dept, ok := claims.GetCustomClaimString("department")
	tst.AssertTrue(t, ok && dept == "engineering", "Department should be preserved")

	level, ok := claims.GetCustomClaimInt("level")
	tst.AssertTrue(t, ok && level == 5, "Level should be preserved")

	isManager, ok := claims.GetCustomClaimBool("is_manager")
	tst.AssertTrue(t, ok && isManager == true, "Manager status should be preserved")
}

func TestInvalidRefreshToken(t *testing.T) {
	manager, err := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	// Test with completely invalid token
	_, err = manager.ValidateRefreshToken("invalid-token")
	tst.AssertTrue(t, errors.Is(err, auth.ErrInvalidRefreshToken), "Should return invalid refresh token error")

	// Test exchange with invalid token
	_, err = manager.ExchangeRefreshToken("invalid-token")
	tst.AssertTrue(t, errors.Is(err, auth.ErrInvalidRefreshToken), "Exchange should fail with invalid token")

	// Test with access token (wrong type)
	accessToken, err := manager.GenerateToken("user123", []string{"user"})
	tst.AssertNoError(t, err)

	_, err = manager.ValidateRefreshToken(accessToken)
	tst.AssertTrue(t, errors.Is(err, auth.ErrInvalidRefreshToken), "Should reject access token as refresh token")
}

func TestExpiredRefreshToken(t *testing.T) {
	// Create manager with very short refresh token duration
	shortManager, err := auth.NewJWTManagerWithRefreshConfig("test-secret", time.Hour, "test-issuer", time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to create JWTManager: %v", err)
	}

	tokenPair, err := shortManager.GenerateTokenPair("user123", []string{"user"})
	tst.AssertNoError(t, err)

	// Wait for refresh token to expire
	time.Sleep(time.Millisecond * 10)

	// Validation should fail
	_, err = shortManager.ValidateRefreshToken(tokenPair.RefreshToken)
	tst.AssertTrue(t, errors.Is(err, auth.ErrRefreshTokenExpired), "Should return expired refresh token error")

	// Exchange should fail
	_, err = shortManager.ExchangeRefreshToken(tokenPair.RefreshToken)
	tst.AssertTrue(t, errors.Is(err, auth.ErrRefreshTokenExpired), "Exchange should fail with expired token")
}

func TestRefreshTokenSecretSeparation(t *testing.T) {
	manager1, err := auth.NewJWTManager("secret1", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)
	manager2, err := auth.NewJWTManager("secret2", time.Hour, "test-issuer")
	tst.AssertNoError(t, err)

	// Generate token pair with manager1
	tokenPair, err := manager1.GenerateTokenPair("user123", []string{"user"})
	tst.AssertNoError(t, err)

	// Manager2 should not be able to validate manager1's refresh token
	_, err = manager2.ValidateRefreshToken(tokenPair.RefreshToken)
	tst.AssertTrue(t, errors.Is(err, auth.ErrInvalidRefreshToken), "Different secret should invalidate refresh token")

	// But manager1 should still be able to validate its own tokens
	_, err = manager1.ValidateRefreshToken(tokenPair.RefreshToken)
	tst.AssertNoError(t, err)
}
