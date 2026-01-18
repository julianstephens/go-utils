package auth

import (
	"errors"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/julianstephens/go-utils/security"
)

var (
	// ErrInvalidToken is returned when a token is malformed or invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired is returned when a token has expired
	ErrTokenExpired = errors.New("token has expired")
	// ErrInvalidClaims is returned when token claims are invalid
	ErrInvalidClaims = errors.New("invalid token claims")
	// ErrInvalidRefreshToken is returned when a refresh token is invalid
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	// ErrRefreshTokenExpired is returned when a refresh token has expired
	ErrRefreshTokenExpired = errors.New("refresh token has expired")
)

const (
	KEY_LENGTH             = 32 // 256 bits for HMAC-SHA256
	ACCESS_SALT            = "go-utils/httputil/auth:access:v1"
	REFRESH_SALT           = "go-utils/httputil/auth:refresh:v1"
	REFRESH_TOKEN_DURATION = time.Hour * 24 * 7 // Default 7 days for refresh tokens
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID       string         `json:"user_id"`
	Username     string         `json:"username,omitempty"`
	Email        string         `json:"email,omitempty"`
	Roles        []string       `json:"roles,omitempty"`
	IssuedAt     time.Time      `json:"iat"`
	CustomClaims map[string]any `json:"custom_claims,omitempty"`
	jwt.RegisteredClaims
}

// RefreshClaims represents the refresh token claims structure
type RefreshClaims struct {
	UserID       string         `json:"user_id"`
	Username     string         `json:"username,omitempty"`
	Email        string         `json:"email,omitempty"`
	Roles        []string       `json:"roles,omitempty"`
	CustomClaims map[string]any `json:"custom_claims,omitempty"`
	TokenID      string         `json:"token_id"` // Unique identifier for this refresh token
	jwt.RegisteredClaims
}

// TokenPair represents both access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // Access token expiration in seconds
}

// KeyPair holds derived keys for access and refresh tokens
type KeyPair struct {
	AccessKey  []byte
	RefreshKey []byte
}

// JWTManager handles JWT token creation and validation
type JWTManager struct {
	secretKey             []byte
	tokenDuration         time.Duration
	issuer                string
	refreshTokenDuration  time.Duration
	refreshTokenSecretKey []byte
}

// NewJWTManager creates a new JWT manager with the given secret key and token duration
func NewJWTManager(secretKey string, tokenDuration time.Duration, issuer string) (*JWTManager, error) {
	keys, err := deriveKeys([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	return &JWTManager{
		secretKey:             keys.AccessKey,
		tokenDuration:         tokenDuration,
		issuer:                issuer,
		refreshTokenDuration:  REFRESH_TOKEN_DURATION,
		refreshTokenSecretKey: keys.RefreshKey,
	}, nil
}

// NewJWTManagerWithRefreshConfig creates a new JWT manager with custom refresh token duration
func NewJWTManagerWithRefreshConfig(
	secretKey string,
	tokenDuration time.Duration,
	issuer string,
	refreshTokenDuration time.Duration,
) (*JWTManager, error) {
	keys, err := deriveKeys([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	return &JWTManager{
		secretKey:             keys.AccessKey,
		tokenDuration:         tokenDuration,
		issuer:                issuer,
		refreshTokenDuration:  refreshTokenDuration,
		refreshTokenSecretKey: keys.RefreshKey,
	}, nil
}

// GenerateTokenPair creates access and refresh token pair
func (j *JWTManager) GenerateTokenPair(userID string, roles []string) (*TokenPair, error) {
	return j.GenerateTokenPairWithClaims(userID, roles, nil)
}

// GenerateTokenPairWithClaims creates access and refresh token pair with custom claims
func (j *JWTManager) GenerateTokenPairWithClaims(
	userID string,
	roles []string,
	customClaims map[string]any,
) (*TokenPair, error) {
	// Generate access token
	accessToken, err := j.GenerateTokenWithClaims(userID, roles, customClaims)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := j.generateRefreshToken(userID, roles, customClaims)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(j.tokenDuration.Seconds()),
	}, nil
}

// GenerateTokenPairWithUserInfo creates access and refresh token pair with user info
func (j *JWTManager) GenerateTokenPairWithUserInfo(userID, username, email string, roles []string) (*TokenPair, error) {
	customClaims := make(map[string]any)
	if username != "" {
		customClaims["username"] = username
	}
	if email != "" {
		customClaims["email"] = email
	}

	if len(customClaims) == 0 {
		customClaims = nil
	}

	return j.GenerateTokenPairWithClaims(userID, roles, customClaims)
}

// generateRefreshToken creates a refresh token with longer expiration
func (j *JWTManager) generateRefreshToken(userID string, roles []string, customClaims map[string]any) (string, error) {
	now := time.Now()
	tokenID := uuid.New().String()

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

	claims := RefreshClaims{
		UserID:       userID,
		Username:     username,
		Email:        email,
		Roles:        roles,
		CustomClaims: customClaims,
		TokenID:      tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   userID,
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.refreshTokenSecretKey)
}

// GenerateToken creates a new JWT token with the provided claims
func (j *JWTManager) GenerateToken(userID string, roles []string) (string, error) {
	return j.GenerateTokenWithClaims(userID, roles, nil)
}

// GenerateTokenWithClaims creates a new JWT token with the provided claims and custom claims
func (j *JWTManager) GenerateTokenWithClaims(
	userID string,
	roles []string,
	customClaims map[string]any,
) (string, error) {
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
	customClaims := make(map[string]any)
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
func (j *JWTManager) GenerateTokenWithUserInfoAndClaims(
	userID, username, email string,
	roles []string,
	customClaims map[string]any,
) (string, error) {
	if customClaims == nil {
		customClaims = make(map[string]any)
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
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
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
		token, parseErr := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
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
		refreshCustomClaims = make(map[string]any)
	}
	if claims.Username != "" {
		refreshCustomClaims["username"] = claims.Username
	}
	if claims.Email != "" {
		refreshCustomClaims["email"] = claims.Email
	}

	return j.GenerateTokenWithClaims(claims.UserID, claims.Roles, refreshCustomClaims)
}

// ValidateRefreshToken validates a refresh token and returns the claims if valid
func (j *JWTManager) ValidateRefreshToken(refreshTokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(refreshTokenString, &RefreshClaims{}, func(token *jwt.Token) (any, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidRefreshToken
		}
		return j.refreshTokenSecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrRefreshTokenExpired
		}
		return nil, ErrInvalidRefreshToken
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidRefreshToken
	}

	return claims, nil
}

// ExchangeRefreshToken validates a refresh token and issues a new token pair
func (j *JWTManager) ExchangeRefreshToken(refreshTokenString string) (*TokenPair, error) {
	// Validate the refresh token
	refreshClaims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Generate new token pair with the same claims
	return j.GenerateTokenPairWithClaims(refreshClaims.UserID, refreshClaims.Roles, refreshClaims.CustomClaims)
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
	return slices.Contains(c.Roles, role)
}

// HasAnyRole checks if the user has any of the specified roles
func (c *Claims) HasAnyRole(roles ...string) bool {
	return slices.ContainsFunc(roles, c.HasRole)
}

// IsExpired checks if the token is expired
func (c *Claims) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(c.ExpiresAt.Time)
}

// Expiration returns the token expiration time and true if present.
// If the token has no expiration claim the returned bool is false.
func (c *Claims) Expiration() (time.Time, bool) {
	if c.ExpiresAt == nil {
		return time.Time{}, false
	}
	return c.ExpiresAt.Time, true
}

// TokenExpiration returns the expiration time for the provided token string.
// It will return the expiration even if the token is expired (by parsing
// without claims validation when needed). If the token is invalid or does
// not contain an expiration, an error is returned.
func (j *JWTManager) TokenExpiration(tokenString string) (time.Time, error) {
	// Try full validation first - this covers valid, non-expired tokens.
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		// If the token is expired, parse without validation to extract claims.
		if !errors.Is(err, ErrTokenExpired) {
			return time.Time{}, err
		}

		token, parseErr := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return j.secretKey, nil
		}, jwt.WithoutClaimsValidation())

		if parseErr != nil {
			return time.Time{}, ErrInvalidToken
		}

		c, ok := token.Claims.(*Claims)
		if !ok {
			return time.Time{}, ErrInvalidClaims
		}

		if c.ExpiresAt == nil {
			return time.Time{}, ErrInvalidClaims
		}
		return c.ExpiresAt.Time, nil
	}

	if claims.ExpiresAt == nil {
		return time.Time{}, ErrInvalidClaims
	}
	return claims.ExpiresAt.Time, nil
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

func deriveKeys(secretKey []byte) (*KeyPair, error) {
	accessKey, refreshKey, err := security.DeriveKeyPair(
		secretKey,
		ACCESS_SALT,
		REFRESH_SALT,
		"access key",
		"refresh key",
		KEY_LENGTH,
	)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		AccessKey:  accessKey,
		RefreshKey: refreshKey,
	}, nil
}
