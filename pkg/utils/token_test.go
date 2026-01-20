package utils

import (
	"os"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	// Setup
	os.Setenv("TOKEN_HOUR_LIFESPAN", "24")
	os.Setenv("API_SECRET", "testsecret123")

	// Test token generation
	userID := uint(1)
	token, err := GenerateToken(userID)

	if err != nil {
		t.Errorf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("GenerateToken returned empty token")
	}

	// Token should be a valid JWT format (three parts separated by dots)
	parts := 0
	for _, c := range token {
		if c == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Errorf("Token should have 3 parts (2 dots), got %d dots", parts)
	}
}

func TestGenerateToken_InvalidLifespan(t *testing.T) {
	// Setup with invalid lifespan
	os.Setenv("TOKEN_HOUR_LIFESPAN", "invalid")
	os.Setenv("API_SECRET", "testsecret123")

	_, err := GenerateToken(1)

	if err == nil {
		t.Error("GenerateToken should fail with invalid TOKEN_HOUR_LIFESPAN")
	}
}

func TestGenerateToken_DifferentUsers(t *testing.T) {
	// Setup
	os.Setenv("TOKEN_HOUR_LIFESPAN", "24")
	os.Setenv("API_SECRET", "testsecret123")

	// Generate tokens for different users
	token1, _ := GenerateToken(1)
	token2, _ := GenerateToken(2)

	if token1 == token2 {
		t.Error("Tokens for different users should be different")
	}
}
