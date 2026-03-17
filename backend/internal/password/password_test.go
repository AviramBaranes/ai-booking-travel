package password

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	pw := "MySecretPass123!"
	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("Hash is empty")
	}
	if hash == pw {
		t.Fatal("Hash should not be equal to password")
	}
}

func TestComparePassword(t *testing.T) {
	pw := "MySecretPass123!"
	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !ComparePassword(hash, pw) {
		t.Error("ComparePassword failed for correct password")
	}

	if ComparePassword(hash, "WrongPass") {
		t.Error("ComparePassword succeeded for wrong password")
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "Valid password",
			password: "StrongPassword1!",
			wantErr:  nil,
		},
		{
			name:     "Too short",
			password: "Short1!",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "No upper",
			password: "lower123!",
			wantErr:  ErrPasswordNoUpper,
		},
		{
			name:     "No lower",
			password: "UPPER123!",
			wantErr:  ErrPasswordNoLower,
		},
		{
			name:     "No number",
			password: "NoNumber!",
			wantErr:  ErrPasswordNoNumber,
		},
		{
			name:     "No symbol",
			password: "NoSymbol1",
			wantErr:  ErrPasswordNoSymbol,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err != tt.wantErr {
					t.Errorf("Expected error %v, got %v", tt.wantErr, err)
				}
			}
		})
	}
}
