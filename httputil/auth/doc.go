// Package auth provides JWT and authentication-related utilities used by
// the httputil subpackages and tests.
//
// The package contains helpers for creating and validating JWT tokens,
// extracting tokens from HTTP headers, and working with claims and roles.
// It is intentionally lightweight and test-friendly.
//
// # Key helpers
//
//   - NewJWTManager(secret []byte, opts ...Option) *JWTManager
//     Creates a token manager for generating and validating tokens.
//
//   - (m *JWTManager) GenerateToken(claims Claims, ttl time.Duration) (string, error)
//     Generate a signed token with a time-to-live.
//
//   - ExtractTokenFromHeader(r *http.Request) (string, error)
//     Extracts the bearer token from an Authorization header.
//
// Example: generate and validate a token
//
//	mgr := auth.NewJWTManager([]byte("my-secret"))
//	token, err := mgr.GenerateToken(auth.Claims{UserID: 42, Roles: []string{"admin"}}, time.Hour)
//	if err != nil { return err }
//
//	claims, err := mgr.ValidateToken(token)
//	if err != nil { return err }
//	fmt.Println("user id:", claims.UserID)
package auth
