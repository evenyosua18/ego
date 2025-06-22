package cryptox

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestHashValue(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantBase64  string
		wantHashHex string
	}{
		{
			name:       "Basic string",
			input:      "test-value",
			wantBase64: base64.URLEncoding.EncodeToString([]byte("test-value")),
		},
		{
			name:       "Empty string",
			input:      "",
			wantBase64: base64.URLEncoding.EncodeToString([]byte("")),
		},
		{
			name:       "Special characters",
			input:      "!@#$_+-=()[]{}",
			wantBase64: base64.URLEncoding.EncodeToString([]byte("!@#$_+-=()[]{}")),
		},
		{
			name:       "Unicode characters",
			input:      "秘密值", // means "secret value" in Chinese
			wantBase64: base64.URLEncoding.EncodeToString([]byte("秘密值")),
		},
		{
			name:       "Long string",
			input:      "a-very-long-secret-value-that-goes-on-and-on-and-might-even-include-💥💣🔥",
			wantBase64: base64.URLEncoding.EncodeToString([]byte("a-very-long-secret-value-that-goes-on-and-on-and-might-even-include-💥💣🔥")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBase64, gotHash := HashValue(tt.input)

			// Compute expected hash on the fly if not provided
			expectedHash := tt.wantHashHex
			if expectedHash == "" {
				sum := sha256.Sum256([]byte(tt.wantBase64))
				expectedHash = base64.URLEncoding.EncodeToString(sum[:])
			}

			if gotBase64 != tt.wantBase64 {
				t.Errorf("base64 mismatch: got %q, want %q", gotBase64, tt.wantBase64)
			}

			if gotHash != expectedHash {
				t.Errorf("hash mismatch: got %q, want %q", gotHash, expectedHash)
			}
		})
	}
}

func TestIsHashValid(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectOK  bool
		useHashOf string // use this as the input to generate stored hash
	}{
		{
			name:      "Valid match",
			value:     "valid-token",
			expectOK:  true,
			useHashOf: "valid-token",
		},
		{
			name:      "Invalid value",
			value:     "wrong-token",
			expectOK:  false,
			useHashOf: "valid-token",
		},
		{
			name:      "Empty string",
			value:     "",
			expectOK:  true,
			useHashOf: "",
		},
		{
			name:      "Unicode value",
			value:     "秘密值",
			expectOK:  true,
			useHashOf: "秘密值",
		},
		{
			name:      "Valid base64-encoded input",
			value:     base64.URLEncoding.EncodeToString([]byte("abc123")),
			expectOK:  true,
			useHashOf: base64.URLEncoding.EncodeToString([]byte("abc123")),
		},
		{
			name:      "Mismatch base64-encoded input",
			value:     base64.URLEncoding.EncodeToString([]byte("xyz456")),
			expectOK:  false,
			useHashOf: base64.URLEncoding.EncodeToString([]byte("abc123")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum := sha256.Sum256([]byte(tt.useHashOf))
			storedHash := base64.URLEncoding.EncodeToString(sum[:])

			isValid := IsHashValid(storedHash, tt.value)
			if isValid != tt.expectOK {
				t.Errorf("IsHashValid(%q) = %v, want %v", tt.value, isValid, tt.expectOK)
			}
		})
	}
}
