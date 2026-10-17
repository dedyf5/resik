// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package uuid

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
)

func TestNewUUIDV7(t *testing.T) {
	id := NewUUIDV7()

	// 1. Verify that the generated UUID is not a Nil UUID
	if id == Nil {
		t.Fatal("NewUUIDV7() returned a Nil UUID")
	}

	// 2. Verify that the version of the generated UUID is correctly set to 7
	version := Version(uuid.UUID(id))
	if version != 7 {
		t.Errorf("NewUUIDV7() version = %d; want 7", version)
	}
}

func TestParseUUIDV7(t *testing.T) {
	validV7Str := NewUUIDV7().String()
	// A valid UUIDv4 string to test version mismatch validation
	validV4Str := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name        string
		input       string
		isErr       bool
		expectedErr error
	}{
		{
			name:        "Valid UUIDv7 string",
			input:       validV7Str,
			isErr:       false,
			expectedErr: nil,
		},
		{
			name:        "Invalid UUID format string",
			input:       "not-a-valid-uuid",
			isErr:       true,
			expectedErr: nil, // Expecting a parse error from uuid package
		},
		{
			name:        "Valid UUID but wrong version (v4)",
			input:       validV4Str,
			isErr:       true,
			expectedErr: ErrInvalidVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUUIDV7(tt.input)

			if tt.isErr {
				if err == nil {
					t.Fatalf("ParseUUIDV7() expected an error, got nil")
				}
				// If a specific error is expected, verify it with errors.Is
				if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
					t.Errorf("ParseUUIDV7() error = %v; want %v", err, tt.expectedErr)
				}
			} else {
				if err != nil {
					t.Fatalf("ParseUUIDV7() unexpected error: %v", err)
				}
				if got == Nil {
					t.Errorf("ParseUUIDV7() got Nil UUID, want valid UUIDv7")
				}
				// Verify the parsed result is actually version 7
				if Version(uuid.UUID(got)) != 7 {
					t.Errorf("ParseUUIDV7() parsed version = %d; want 7", Version(uuid.UUID(got)))
				}
			}
		})
	}
}

func TestUUIDV7_String32(t *testing.T) {
	id := NewUUIDV7()
	str32 := id.String32()

	// 1. Verify the exact length (standard UUID without dashes is 32 characters)
	expectedLength := 32
	if len(str32) != expectedLength {
		t.Errorf("UUIDV7.String32() length = %d; want %d", len(str32), expectedLength)
	}

	// 2. Verify that there are no dash characters ('-') present
	if strings.Contains(str32, "-") {
		t.Errorf("UUIDV7.String32() contains dashes: %q", str32)
	}

	// 3. Verify consistency by comparing with standard string representation stripped of dashes
	standardStr := id.uuid().String()
	expectedStr32 := strings.ReplaceAll(standardStr, "-", "")
	if str32 != expectedStr32 {
		t.Errorf("UUIDV7.String32() = %q; want %q", str32, expectedStr32)
	}
}

func TestUUIDV7_Equal(t *testing.T) {
	id := NewUUIDV7()
	idSame := id
	idDifferent := NewUUIDV7()

	tests := []struct {
		name     string
		input    UUIDV7
		expected bool
	}{
		{
			name:     "Equal",
			input:    idSame,
			expected: true,
		},
		{
			name:     "Different",
			input:    idDifferent,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := id.Equal(tt.input.String())
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestUUIDV7_MarshalJSON(t *testing.T) {
	id := NewUUIDV7()

	// 1. Test direct MarshalJSON method call
	data, err := id.MarshalJSON()
	if err != nil {
		t.Fatalf("UUIDV7.MarshalJSON() unexpected error: %v", err)
	}

	// Verify the exact length (32 hex characters + 2 double quotes = 34 bytes)
	expectedLength := 34
	if len(data) != expectedLength {
		t.Errorf("UUIDV7.MarshalJSON() length = %d; want %d", len(data), expectedLength)
	}

	// Verify that the output starts and ends with double quotes
	if data[0] != '"' || data[len(data)-1] != '"' {
		t.Errorf("UUIDV7.MarshalJSON() result %q is missing proper quotation marks", string(data))
	}

	// 2. Test integration with standard encoding/json package
	jsonData, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("json.Marshal() unexpected error: %v", err)
	}

	// Verify that json.Marshal produces the exact same result as custom MarshalJSON
	if string(jsonData) != string(data) {
		t.Errorf("json.Marshal() = %q; want %q", string(jsonData), string(data))
	}
}

func TestUUIDV7_UnmarshalJSON(t *testing.T) {
	original := NewUUIDV7()
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal UUIDv7 for test setup: %v", err)
	}

	t.Run("Valid JSON unmarshal", func(t *testing.T) {
		var target UUIDV7
		err := target.UnmarshalJSON(jsonData)
		if err != nil {
			t.Fatalf("UUIDV7.UnmarshalJSON() unexpected error: %v", err)
		}

		// Verify that the unmarshaled value matches the original
		if target != original {
			t.Errorf("UUIDV7.UnmarshalJSON() result = %v; want %v", target, original)
		}
	})

	t.Run("Nil data input", func(t *testing.T) {
		var target UUIDV7
		err := target.UnmarshalJSON(nil)
		if err != nil {
			t.Errorf("UUIDV7.UnmarshalJSON(nil) expected no error, got %v", err)
		}
	})

	t.Run("Invalid JSON format", func(t *testing.T) {
		var target UUIDV7
		// Invalid JSON syntax (e.g., malformed object instead of a JSON string)
		invalidJSON := []byte(`{invalid-json}`)
		err := target.UnmarshalJSON(invalidJSON)
		if err == nil {
			t.Error("UUIDV7.UnmarshalJSON() expected error for invalid JSON, got nil")
		}
	})

	t.Run("Invalid UUID string inside JSON", func(t *testing.T) {
		var target UUIDV7
		// Valid JSON string syntax, but the content is not a valid UUID
		invalidUUIDJSON := []byte(`"not-a-valid-uuid-string"`)
		err := target.UnmarshalJSON(invalidUUIDJSON)
		if err == nil {
			t.Error("UUIDV7.UnmarshalJSON() expected error for invalid UUID string, got nil")
		}
	})
}

func TestUUIDV7_IsEmpty(t *testing.T) {
	id := NewUUIDV7()
	if id.IsEmpty() {
		t.Error("UUIDV7.IsEmpty() returned true for a non-empty UUID")
	}

	emptyUUID := Nil
	if !emptyUUID.IsEmpty() {
		t.Error("UUIDV7.IsEmpty() returned false for an empty UUID")
	}
}

func TestUUIDV7_Value(t *testing.T) {
	id := NewUUIDV7()

	// Call Value method
	val, err := id.Value()
	if err != nil {
		t.Fatalf("UUIDV7.Value() unexpected error: %v", err)
	}

	// Verify that the returned value is a string type
	strVal, ok := val.(string)
	if !ok {
		t.Fatalf("UUIDV7.Value() returned type %T; want string", val)
	}

	// Verify that the string value matches the standard string representation
	expectedStr := id.uuid().String()
	if strVal != expectedStr {
		t.Errorf("UUIDV7.Value() = %q; want %q", strVal, expectedStr)
	}
}

func TestUUIDV7_Scan(t *testing.T) {
	validUUIDStr := NewUUIDV7().String()

	tests := []struct {
		name        string
		input       any
		expectError bool
		expectedErr error
	}{
		{
			name:        "Scan nil value",
			input:       nil,
			expectError: false,
			expectedErr: nil,
		},
		{
			name:        "Scan valid byte slice",
			input:       []byte(validUUIDStr),
			expectError: false,
			expectedErr: nil,
		},
		{
			name:        "Scan invalid type (string instead of []byte)",
			input:       validUUIDStr,
			expectError: true,
			expectedErr: ErrInvalidValue,
		},
		{
			name:        "Scan invalid type (int)",
			input:       12345,
			expectError: true,
			expectedErr: ErrInvalidValue,
		},
		{
			name:        "Scan invalid UUID format in byte slice",
			input:       []byte("not-a-valid-uuid"),
			expectError: true,
			expectedErr: nil, // Expecting parse error from uuid package
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target UUIDV7
			err := target.Scan(tt.input)

			if tt.expectError {
				if err == nil {
					t.Fatalf("UUIDV7.Scan() expected an error, got nil")
				}
				if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
					t.Errorf("UUIDV7.Scan() error = %v; want %v", err, tt.expectedErr)
				}
			} else {
				if err != nil {
					t.Fatalf("UUIDV7.Scan() unexpected error: %v", err)
				}
				// If input was a valid byte slice, verify target was populated correctly
				if tt.input != nil {
					bytesVal := tt.input.([]byte)
					expectedParsed, _ := uuid.Parse(string(bytesVal))
					if target != UUIDV7(expectedParsed) {
						t.Errorf("UUIDV7.Scan() result = %v; want %v", target, UUIDV7(expectedParsed))
					}
				}
			}
		})
	}
}
