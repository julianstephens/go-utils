package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/julianstephens/go-utils/httputil/response"
)

// RefreshTokenHandler creates an HTTP handler for token exchange
// It expects a JSON payload with a "refresh_token" field and returns a new TokenPair
func RefreshTokenHandler(manager *JWTManager) http.HandlerFunc {
	responder := response.NewEmpty()
	
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			responder.ErrorWithStatus(w, r, http.StatusMethodNotAllowed, errors.New("only POST method allowed"))
			return
		}

		var request struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			responder.BadRequest(w, r, "invalid JSON payload")
			return
		}

		if request.RefreshToken == "" {
			responder.BadRequest(w, r, "refresh_token is required")
			return
		}

		// Exchange the refresh token for a new token pair
		tokenPair, err := manager.ExchangeRefreshToken(request.RefreshToken)
		if err != nil {
			if errors.Is(err, ErrRefreshTokenExpired) {
				responder.Unauthorized(w, r, "refresh token expired")
				return
			}
			if errors.Is(err, ErrInvalidRefreshToken) {
				responder.Unauthorized(w, r, "invalid refresh token")
				return
			}
			responder.InternalServerError(w, r, "failed to exchange token")
			return
		}

		responder.OK(w, r, tokenPair)
	}
}

// AuthenticationHandler creates an HTTP handler for user authentication
// This is an example handler that shows how to issue token pairs during login
func AuthenticationHandler(manager *JWTManager, authenticateUser func(username, password string) (*UserInfo, error)) http.HandlerFunc {
	responder := response.NewEmpty()

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			responder.ErrorWithStatus(w, r, http.StatusMethodNotAllowed, errors.New("only POST method allowed"))
			return
		}

		var request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			responder.BadRequest(w, r, "invalid JSON payload")
			return
		}

		if request.Username == "" || request.Password == "" {
			responder.BadRequest(w, r, "username and password are required")
			return
		}

		// Authenticate user (implementation depends on your user store)
		userInfo, err := authenticateUser(request.Username, request.Password)
		if err != nil {
			responder.Unauthorized(w, r, "invalid credentials")
			return
		}

		// Generate token pair for authenticated user
		tokenPair, err := manager.GenerateTokenPairWithUserInfo(
			userInfo.UserID,
			userInfo.Username,
			userInfo.Email,
			userInfo.Roles,
		)
		if err != nil {
			responder.InternalServerError(w, r, "failed to generate tokens")
			return
		}

		// Return token pair
		responder.OK(w, r, tokenPair)
	}
}

// UserInfo represents user information for authentication
type UserInfo struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

// SetRefreshTokenCookie sets a refresh token as an httpOnly cookie
// This is a security best practice for web applications
func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string, maxAge time.Duration, secure bool) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: true,
		Secure:   secure, // Should be true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

// GetRefreshTokenFromCookie extracts refresh token from httpOnly cookie
func GetRefreshTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// ClearRefreshTokenCookie clears the refresh token cookie (for logout)
func ClearRefreshTokenCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}