package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken is returned when a token is malformed or invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired is returned when a token has expired
	ErrTokenExpired = errors.New("token has expired")
	// ErrInvalidClaims is returned when token claims are invalid
	ErrInvalidClaims = errors.New("invalid token claims")
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID       string                 `json:"user_id"`
	Username     string                 `json:"username,omitempty"`
	Email        string                 `json:"email,omitempty"`
	Roles        []string               `json:"roles,omitempty"`
	IssuedAt     time.Time              `json:"iat"`
	CustomClaims map[string]interface{} `json:"custom_claims,omitempty"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token creation and validation
type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
	issuer        string
}

// NewJWTManager creates a new JWT manager with the given secret key and token duration
func NewJWTManager(secretKey string, tokenDuration time.Duration, issuer string) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
		issuer:        issuer,
	}
}

// GenerateToken creates a new JWT token with the provided claims
func (j *JWTManager) GenerateToken(userID string, roles []string) (string, error) {
	return j.GenerateTokenWithClaims(userID, roles, nil)
}

// GenerateTokenWithClaims creates a new JWT token with the provided claims and custom claims
func (j *JWTManager) GenerateTokenWithClaims(userID string, roles []string, customClaims map[string]interface{}) (string, error) {
	now := time.Now()
	
	// Extract username and email from custom claims if provided
	var username, email string
	if customClaims != nil {
		if u, ok := customClaims["username"].(string); ok {
			username = u
		}
		if e, ok := customClaims["email"].(string); ok {
			email = e
		}
	}
	
	claims := Claims{
		UserID:       userID,
		Username:     username,
		Email:        email,
		Roles:        roles,
		IssuedAt:     now,
		CustomClaims: customClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// GenerateTokenWithUserInfo creates a new JWT token with user ID, username, email, and roles
// This is a convenience method for backward compatibility
func (j *JWTManager) GenerateTokenWithUserInfo(userID, username, email string, roles []string) (string, error) {
	customClaims := make(map[string]interface{})
	if username != "" {
		customClaims["username"] = username
	}
	if email != "" {
		customClaims["email"] = email
	}
	
	if len(customClaims) == 0 {
		customClaims = nil
	}
	
	return j.GenerateTokenWithClaims(userID, roles, customClaims)
}

// GenerateTokenWithUserInfoAndClaims creates a new JWT token with user info and additional custom claims
// This is a convenience method for backward compatibility
func (j *JWTManager) GenerateTokenWithUserInfoAndClaims(userID, username, email string, roles []string, customClaims map[string]interface{}) (string, error) {
	if customClaims == nil {
		customClaims = make(map[string]interface{})
	}
	
	if username != "" {
		customClaims["username"] = username
	}
	if email != "" {
		customClaims["email"] = email
	}
	
	return j.GenerateTokenWithClaims(userID, roles, customClaims)
}

// ValidateToken validates a JWT token and returns the claims if valid
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// RefreshToken generates a new token with updated expiration time for valid existing token
func (j *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		// Allow refresh even if token is expired
		if !errors.Is(err, ErrTokenExpired) {
			return "", err
		}
		// Parse expired token to get claims
		token, parseErr := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return j.secretKey, nil
		}, jwt.WithoutClaimsValidation())

		if parseErr != nil {
			return "", ErrInvalidToken
		}

		var ok bool
		claims, ok = token.Claims.(*Claims)
		if !ok {
			return "", ErrInvalidClaims
		}
	}

	// Generate new token with same claims but updated timestamps
	// Preserve username and email in custom claims if they exist
	refreshCustomClaims := claims.CustomClaims
	if refreshCustomClaims == nil {
		refreshCustomClaims = make(map[string]interface{})
	}
	if claims.Username != "" {
		refreshCustomClaims["username"] = claims.Username
	}
	if claims.Email != "" {
		refreshCustomClaims["email"] = claims.Email
	}
	
	return j.GenerateTokenWithClaims(claims.UserID, claims.Roles, refreshCustomClaims)
}

// ExtractTokenFromHeader extracts JWT token from Authorization header
// Expected format: "Bearer <token>"
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("invalid authorization header format")
	}

	return authHeader[len(bearerPrefix):], nil
}

// HasRole checks if the user has a specific role
func (c *Claims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole checks if the user has any of the specified roles
func (c *Claims) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if c.HasRole(role) {
			return true
		}
	}
	return false
}

// IsExpired checks if the token is expired
func (c *Claims) IsExpired() bool {
	if c.RegisteredClaims.ExpiresAt == nil {
		return false
	}
	return time.Now().After(c.RegisteredClaims.ExpiresAt.Time)
}

// GetCustomClaim retrieves a custom claim value by key
func (c *Claims) GetCustomClaim(key string) (interface{}, bool) {
	if c.CustomClaims == nil {
		return nil, false
	}
	value, exists := c.CustomClaims[key]
	return value, exists
}

// GetCustomClaimString retrieves a custom claim as a string
func (c *Claims) GetCustomClaimString(key string) (string, bool) {
	value, exists := c.GetCustomClaim(key)
	if !exists {
		return "", false
	}
	str, ok := value.(string)
	return str, ok
}

// GetCustomClaimInt retrieves a custom claim as an int
func (c *Claims) GetCustomClaimInt(key string) (int, bool) {
	value, exists := c.GetCustomClaim(key)
	if !exists {
		return 0, false
	}

	// Handle different numeric types that JSON might unmarshal to
	switch v := value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case int64:
		return int(v), true
	default:
		return 0, false
	}
}

// GetCustomClaimBool retrieves a custom claim as a boolean
func (c *Claims) GetCustomClaimBool(key string) (bool, bool) {
	value, exists := c.GetCustomClaim(key)
	if !exists {
		return false, false
	}
	boolean, ok := value.(bool)
	return boolean, ok
}

// HasCustomClaim checks if a custom claim exists
func (c *Claims) HasCustomClaim(key string) bool {
	if c.CustomClaims == nil {
		return false
	}
	_, exists := c.CustomClaims[key]
	return exists
}

// SetCustomClaim sets a custom claim (useful for building claims before token generation)
func (c *Claims) SetCustomClaim(key string, value interface{}) {
	if c.CustomClaims == nil {
		c.CustomClaims = make(map[string]interface{})
	}
	c.CustomClaims[key] = value
}

// DeleteCustomClaim removes a custom claim
func (c *Claims) DeleteCustomClaim(key string) {
	if c.CustomClaims != nil {
		delete(c.CustomClaims, key)
	}
}
