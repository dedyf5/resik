// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package hash

import (
	"crypto/rand"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArgon2Hasher(t *testing.T) {
	tests := []struct {
		name     string
		input    *Argon2Config
		expected *Argon2Config
	}{
		{
			name: "Default values when config is zero",
			input: &Argon2Config{
				Memory:      0,
				Iterations:  0,
				Parallelism: 0,
				SaltLength:  0,
				KeyLength:   0,
			},
			expected: &Argon2Config{
				Memory:      8 * 1024,
				Iterations:  3,
				Parallelism: 1,
				SaltLength:  16,
				KeyLength:   32,
			},
		},
		{
			name: "Values below minimum should be sanitized",
			input: &Argon2Config{
				Memory:      1024,
				Iterations:  0,
				Parallelism: 0,
				SaltLength:  0,
				KeyLength:   5,
			},
			expected: &Argon2Config{
				Memory:      8 * 1024,
				Iterations:  3,
				Parallelism: 1,
				SaltLength:  16,
				KeyLength:   5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasherInterface := NewArgon2Hasher(tt.input)

			h, ok := hasherInterface.(*argon2Hasher)

			assert.True(t, ok)
			assert.Equal(t, tt.expected.Memory, h.config.Memory)
			assert.Equal(t, tt.expected.Iterations, h.config.Iterations)
			assert.Equal(t, tt.expected.Parallelism, h.config.Parallelism)
			assert.Equal(t, tt.expected.SaltLength, h.config.SaltLength)
			assert.Equal(t, tt.expected.KeyLength, h.config.KeyLength)
		})
	}
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated randomness failure")
}

func TestArgon2Hasher_Hash(t *testing.T) {
	config := &Argon2Config{
		Memory: 8 * 1024, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32,
	}

	t.Run("Fail - When randomness source fails", func(t *testing.T) {
		hasher := &argon2Hasher{config: config, source: &errorReader{}}
		encodedHash, err := hasher.Hash("any_password")
		assert.Error(t, err)
		assert.Empty(t, encodedHash)
	})

	t.Run("Success - Should generate valid encoded hash", func(t *testing.T) {
		hasher := &argon2Hasher{config: config, source: rand.Reader}
		encodedHash, err := hasher.Hash("secret")
		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(encodedHash, "$argon2id$"))
	})

	t.Run("Success - Unique hashes for same password", func(t *testing.T) {
		hasher := &argon2Hasher{config: config, source: rand.Reader}
		h1, _ := hasher.Hash("password")
		h2, _ := hasher.Hash("password")
		assert.NotEqual(t, h1, h2)
	})
}

func TestArgon2Hasher_Compare(t *testing.T) {
	hasher := NewArgon2Hasher(&Argon2Config{
		Memory: 8 * 1024, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32,
	})
	validHash, _ := hasher.Hash("correct_password")

	tests := []struct {
		name          string
		password      string
		encodedHash   string
		expectedMatch bool
		expectedError error
	}{
		{
			name:          "Fail - Invalid format length",
			encodedHash:   "$short$hash",
			expectedError: ErrInvalidHash,
		},
		{
			name:          "Fail - Wrong algorithm",
			encodedHash:   "$bcrypt$v=19$m=8192,t=1,p=1$salt$hash",
			expectedError: ErrInvalidHash,
		},
		{
			name:          "Fail - Sscanf version malformed",
			encodedHash:   "$argon2id$v=abc$m=8192,t=1,p=1$salt$hash",
			expectedError: errors.New("sscanf error"),
		},
		{
			name:          "Fail - Incompatible version",
			encodedHash:   "$argon2id$v=10$m=8192,t=1,p=1$salt$hash",
			expectedError: ErrIncompatibleVersion,
		},
		{
			name:          "Fail - Sscanf params malformed",
			encodedHash:   "$argon2id$v=19$m=abc,t=1,p=1$salt$hash",
			expectedError: errors.New("sscanf error"),
		},
		{
			name:          "Fail - Invalid base64 salt",
			encodedHash:   "$argon2id$v=19$m=8192,t=1,p=1$invalid_salt!$hash",
			expectedError: errors.New("base64 error"),
		},
		{
			name:          "Fail - Invalid base64 hash",
			encodedHash:   "$argon2id$v=19$m=8192,t=1,p=1$c2FsdA$hash%invalid",
			expectedError: errors.New("base64 error"),
		},
		{
			name:          "Success - Match",
			password:      "correct_password",
			encodedHash:   validHash,
			expectedMatch: true,
		},
		{
			name:          "Fail - Mismatch password",
			password:      "wrong_password",
			encodedHash:   validHash,
			expectedMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := hasher.Compare(tt.password, tt.encodedHash)
			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedMatch, match)
		})
	}
}
