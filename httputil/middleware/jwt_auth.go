package middleware

import (
	"context"
	"net/http"

	"github.com/julianstephens/go-utils/httputil/auth"
	"github.com/julianstephens/go-utils/httputil/response"
)

// AuthClaimsKey is the context key used to store auth.Claims
const AuthClaimsKey contextKey = "auth_claims"

// JWTAuth creates middleware that validates JWT tokens from the
// Authorization header using the provided JWTManager. On success,
// the claims are stored in the request context under AuthClaimsKey.
// If validation fails, a JSON 401 Unauthorized response is returned.
func JWTAuth(manager *auth.JWTManager) func(http.Handler) http.Handler {
	responder := response.NewEmpty()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			tokenString, err := auth.ExtractTokenFromHeader(authHeader)
			if err != nil {
				responder.Unauthorized(w, r, "authorization header required")
				return
			}

			claims, err := manager.ValidateToken(tokenString)
			if err != nil {
				responder.Unauthorized(w, r, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), AuthClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims retrieves auth.Claims from context, if present.
func GetClaims(ctx context.Context) (*auth.Claims, bool) {
	if v := ctx.Value(AuthClaimsKey); v != nil {
		if c, ok := v.(*auth.Claims); ok {
			return c, true
		}
	}
	return nil, false
}

// RequireRoles returns middleware that ensures the authenticated user has at least
// one of the specified roles. If the user is not authenticated, a 401 is returned.
// If authenticated but lacking roles, a 403 Forbidden JSON response is returned.
func RequireRoles(manager *auth.JWTManager, roles ...string) func(http.Handler) http.Handler {
	responder := response.NewEmpty()
	// Use the basic JWTAuth to validate and inject claims first
	authMW := JWTAuth(manager)

	return func(next http.Handler) http.Handler {
		// Compose middleware: first JWTAuth, then role check
		return authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r.Context())
			if !ok {
				responder.Unauthorized(w, r, "authorization header required")
				return
			}

			// If no roles required, allow through
			if len(roles) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			if !claims.HasAnyRole(roles...) {
				responder.Forbidden(w, r, "insufficient role")
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}
