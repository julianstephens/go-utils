package main

import (
	"fmt"
	"log"
	"time"

	"github.com/julianstephens/go-utils/httputil/auth"
)

func main() {
	// Create JWT manager
	manager := auth.NewJWTManager("my-secret-key", time.Hour*24, "my-app")

	// Standard claims
	userID := "user123"
	username := "john_doe"
	email := "john@example.com"
	roles := []string{"user", "admin"}

	// Custom claims - can be any JSON-serializable data
	customClaims := map[string]interface{}{
		"department":  "engineering",
		"level":       5,
		"is_manager":  true,
		"permissions": []string{"read", "write", "delete"},
		"metadata": map[string]string{
			"region": "us-west",
			"tenant": "acme-corp",
		},
		"last_login": time.Now().Unix(),
	}

	// Generate token with custom claims
	token, err := manager.GenerateTokenWithClaims(userID, username, email, roles, customClaims)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	fmt.Printf("Generated JWT Token:\n%s\n\n", token)

	// Validate token
	claims, err := manager.ValidateToken(token)
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}

	fmt.Printf("Standard Claims:\n")
	fmt.Printf("  User ID: %s\n", claims.UserID)
	fmt.Printf("  Username: %s\n", claims.Username)
	fmt.Printf("  Email: %s\n", claims.Email)
	fmt.Printf("  Roles: %v\n", claims.Roles)

	fmt.Printf("\nCustom Claims:\n")

	// Retrieve custom claims using type-safe methods
	if dept, ok := claims.GetCustomClaimString("department"); ok {
		fmt.Printf("  Department: %s\n", dept)
	}

	if level, ok := claims.GetCustomClaimInt("level"); ok {
		fmt.Printf("  Level: %d\n", level)
	}

	if isManager, ok := claims.GetCustomClaimBool("is_manager"); ok {
		fmt.Printf("  Is Manager: %t\n", isManager)
	}

	// Retrieve complex custom claims
	if metadata, ok := claims.GetCustomClaim("metadata"); ok {
		fmt.Printf("  Metadata: %v\n", metadata)
	}

	if lastLogin, ok := claims.GetCustomClaimInt("last_login"); ok {
		loginTime := time.Unix(int64(lastLogin), 0)
		fmt.Printf("  Last Login: %s\n", loginTime.Format("2006-01-02 15:04:05"))
	}

	// Check role-based access
	fmt.Printf("\nRole-based Access:\n")
	fmt.Printf("  Has 'admin' role: %t\n", claims.HasRole("admin"))
	fmt.Printf("  Has any of ['user', 'guest']: %t\n", claims.HasAnyRole("user", "guest"))

	// Demonstrate token refresh (preserves custom claims)
	refreshedToken, err := manager.RefreshToken(token)
	if err != nil {
		log.Fatalf("Failed to refresh token: %v", err)
	}

	refreshedClaims, err := manager.ValidateToken(refreshedToken)
	if err != nil {
		log.Fatalf("Failed to validate refreshed token: %v", err)
	}

	fmt.Printf("\nRefreshed Token Custom Claims Preserved:\n")
	if dept, ok := refreshedClaims.GetCustomClaimString("department"); ok {
		fmt.Printf("  Department in refreshed token: %s\n", dept)
	}
}
