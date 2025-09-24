package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julianstephens/go-utils/httputil/auth"
)

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	manager := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")
	token, err := manager.GenerateToken("user1", []string{"user"})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	handler := JWTAuth(manager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ensure claims are present in context
		if _, ok := GetClaims(r.Context()); !ok {
			t.Error("expected claims in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}

func TestJWTAuthMiddleware_Unauthorized(t *testing.T) {
	manager := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")

	handler := JWTAuth(manager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// missing header
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for missing header, got %d", rr.Code)
	}

	// invalid header
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic abc")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for invalid header, got %d", rr.Code)
	}
}

func TestRequireRoles_AllowedAndForbidden(t *testing.T) {
	manager := auth.NewJWTManager("test-secret", time.Hour, "test-issuer")

	// token with admin role
	adminToken, err := manager.GenerateToken("admin1", []string{"admin", "user"})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// token with only user role
	userToken, err := manager.GenerateToken("user1", []string{"user"})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	handler := RequireRoles(manager, "admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// admin should be allowed
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 for admin, got %d", rr.Code)
	}

	// user should be forbidden
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status 403 for user, got %d", rr.Code)
	}
}
