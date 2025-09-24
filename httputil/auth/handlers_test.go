package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julianstephens/go-utils/httputil/auth"
	tst "github.com/julianstephens/go-utils/tests"
)

func TestRefreshTokenHandler(t *testing.T) {
	manager := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	handler := auth.RefreshTokenHandler(manager)

	// Generate initial token pair
	tokenPair, err := manager.GenerateTokenPair("user123", []string{"user"})
	tst.AssertNoError(t, err)

	// Test successful token exchange
	t.Run("successful token exchange", func(t *testing.T) {
		payload := map[string]string{
			"refresh_token": tokenPair.RefreshToken,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		tst.AssertTrue(t, w.Code == http.StatusOK, "Expected 200 OK")

		var response auth.TokenPair
		err := json.NewDecoder(w.Body).Decode(&response)
		tst.AssertNoError(t, err)
		tst.AssertTrue(t, response.AccessToken != "", "Should return new access token")
		tst.AssertTrue(t, response.RefreshToken != "", "Should return new refresh token")
		tst.AssertTrue(t, response.AccessToken != tokenPair.AccessToken, "New access token should be different")
		tst.AssertTrue(t, response.RefreshToken != tokenPair.RefreshToken, "New refresh token should be different")
	})

	// Test with invalid refresh token
	t.Run("invalid refresh token", func(t *testing.T) {
		payload := map[string]string{
			"refresh_token": "invalid-token",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		tst.AssertTrue(t, w.Code == http.StatusUnauthorized, "Expected 401 Unauthorized")
	})

	// Test with missing refresh token
	t.Run("missing refresh token", func(t *testing.T) {
		payload := map[string]string{}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		tst.AssertTrue(t, w.Code == http.StatusBadRequest, "Expected 400 Bad Request")
	})

	// Test with invalid JSON
	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		tst.AssertTrue(t, w.Code == http.StatusBadRequest, "Expected 400 Bad Request")
	})

	// Test with wrong HTTP method
	t.Run("wrong HTTP method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/refresh", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		tst.AssertTrue(t, w.Code == http.StatusMethodNotAllowed, "Expected 405 Method Not Allowed")
	})
}

func TestAuthenticationHandler(t *testing.T) {
	manager := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")

	// Mock authentication function
	authenticateUser := func(username, password string) (*auth.UserInfo, error) {
		if username == "testuser" && password == "testpass" {
			return &auth.UserInfo{
				UserID:   "user123",
				Username: "testuser",
				Email:    "test@example.com",
				Roles:    []string{"user", "admin"},
			}, nil
		}
		return nil, auth.ErrInvalidToken // Reusing existing error for simplicity
	}

	handler := auth.AuthenticationHandler(manager, authenticateUser)

	// Test successful authentication
	t.Run("successful authentication", func(t *testing.T) {
		payload := map[string]string{
			"username": "testuser",
			"password": "testpass",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		tst.AssertTrue(t, w.Code == http.StatusOK, "Expected 200 OK")

		var response auth.TokenPair
		err := json.NewDecoder(w.Body).Decode(&response)
		tst.AssertNoError(t, err)
		tst.AssertTrue(t, response.AccessToken != "", "Should return access token")
		tst.AssertTrue(t, response.RefreshToken != "", "Should return refresh token")
		tst.AssertTrue(t, response.TokenType == "Bearer", "Token type should be Bearer")

		// Validate the access token contains correct user info
		claims, err := manager.ValidateToken(response.AccessToken)
		tst.AssertNoError(t, err)
		tst.AssertTrue(t, claims.UserID == "user123", "UserID should match")
		tst.AssertTrue(t, claims.Username == "testuser", "Username should match")
		tst.AssertTrue(t, claims.Email == "test@example.com", "Email should match")
		tst.AssertTrue(t, len(claims.Roles) == 2, "Should have correct number of roles")
	})

	// Test invalid credentials
	t.Run("invalid credentials", func(t *testing.T) {
		payload := map[string]string{
			"username": "testuser",
			"password": "wrongpass",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		tst.AssertTrue(t, w.Code == http.StatusUnauthorized, "Expected 401 Unauthorized")
	})

	// Test missing credentials
	t.Run("missing credentials", func(t *testing.T) {
		payload := map[string]string{
			"username": "testuser",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		tst.AssertTrue(t, w.Code == http.StatusBadRequest, "Expected 400 Bad Request")
	})
}

func TestRefreshTokenCookieHelpers(t *testing.T) {
	// Test setting refresh token cookie
	t.Run("set refresh token cookie", func(t *testing.T) {
		w := httptest.NewRecorder()
		refreshToken := "test-refresh-token"
		maxAge := time.Hour * 24

		auth.SetRefreshTokenCookie(w, refreshToken, maxAge, true)

		cookies := w.Result().Cookies()
		tst.AssertTrue(t, len(cookies) == 1, "Should set one cookie")

		cookie := cookies[0]
		tst.AssertTrue(t, cookie.Name == "refresh_token", "Cookie name should be refresh_token")
		tst.AssertTrue(t, cookie.Value == refreshToken, "Cookie value should match")
		tst.AssertTrue(t, cookie.HttpOnly == true, "Cookie should be httpOnly")
		tst.AssertTrue(t, cookie.Secure == true, "Cookie should be secure")
		tst.AssertTrue(t, cookie.SameSite == http.SameSiteStrictMode, "Cookie should be SameSite strict")
		tst.AssertTrue(t, cookie.MaxAge == int(maxAge.Seconds()), "Cookie MaxAge should match")
	})

	// Test getting refresh token from cookie
	t.Run("get refresh token from cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: "test-refresh-token",
		})

		token, err := auth.GetRefreshTokenFromCookie(req)
		tst.AssertNoError(t, err)
		tst.AssertTrue(t, token == "test-refresh-token", "Should return correct token value")
	})

	// Test getting refresh token when cookie doesn't exist
	t.Run("get refresh token no cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		token, err := auth.GetRefreshTokenFromCookie(req)
		tst.AssertTrue(t, err != nil, "Should return error when cookie doesn't exist")
		tst.AssertTrue(t, token == "", "Should return empty token")
	})

	// Test clearing refresh token cookie
	t.Run("clear refresh token cookie", func(t *testing.T) {
		w := httptest.NewRecorder()

		auth.ClearRefreshTokenCookie(w)

		cookies := w.Result().Cookies()
		tst.AssertTrue(t, len(cookies) == 1, "Should set one cookie")

		cookie := cookies[0]
		tst.AssertTrue(t, cookie.Name == "refresh_token", "Cookie name should be refresh_token")
		tst.AssertTrue(t, cookie.Value == "", "Cookie value should be empty")
		tst.AssertTrue(t, cookie.MaxAge == -1, "Cookie MaxAge should be -1 to delete")
	})
}
