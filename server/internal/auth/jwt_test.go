package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateJWTWithExpiration(userID int, exp time.Time) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestGenerateJWT(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	tests := []struct {
		name   string
		userID int
	}{
		{"valid user ID", 1},
		{"another user ID", 42},
		{"large user ID", 999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWT(tt.userID)
			if err != nil {
				t.Errorf("GenerateJWT() error = %v, want nil", err)
				return
			}
			if token == "" {
				t.Error("GenerateJWT() returned empty token")
				return
			}

			// Verify the token can be validated and contains correct userID
			gotUserID, err := ValidateJWT(token)
			if err != nil {
				t.Errorf("ValidateJWT() error = %v, want nil", err)
				return
			}
			if gotUserID != tt.userID {
				t.Errorf("ValidateJWT() = %v, want %v", gotUserID, tt.userID)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	// Generate valid token for tests
	validToken, _ := GenerateJWT(123)

	// Generate expired token
	expiredToken, _ := generateJWTWithExpiration(123, time.Now().Add(-1*time.Hour))

	// Generate token with different secret
	t.Setenv("JWT_SECRET", "different-secret")
	wrongSecretToken, _ := GenerateJWT(123)
	t.Setenv("JWT_SECRET", "test-secret")

	tests := []struct {
		name       string
		token      string
		wantUserID int
		wantErr    bool
	}{
		{
			name:       "valid token",
			token:      validToken,
			wantUserID: 123,
			wantErr:    false,
		},
		{
			name:       "expired token",
			token:      expiredToken,
			wantUserID: 0,
			wantErr:    true,
		},
		{
			name:       "wrong secret token",
			token:      wrongSecretToken,
			wantUserID: 0,
			wantErr:    true,
		},
		{
			name:       "malformed token",
			token:      "not.a.valid.token",
			wantUserID: 0,
			wantErr:    true,
		},
		{
			name:       "empty token",
			token:      "",
			wantUserID: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
