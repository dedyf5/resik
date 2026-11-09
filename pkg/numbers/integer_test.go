// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package numbers

import (
	"errors"
	"testing"
)

func TestSaveConvert_Success(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int16
	}{
		{
			name:  "normal positive value within range",
			input: 150,
			want:  150,
		},
		{
			name:  "zero value",
			input: 0,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeConvert[int16](tt.input)
			if err != nil {
				t.Errorf("SafeConvert() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Errorf("SafeConvert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveConvert_Overflow(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{
			name:  "exceeds int8 max capacity",
			input: 200,
		},
		{
			name:  "below int8 min capacity",
			input: -200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SafeConvert[int8](tt.input)
			if !errors.Is(err, ErrOverflow) {
				t.Errorf("SafeConvert() error = %v, want %v", err, ErrOverflow)
			}
		})
	}
}

func TestSaveConvert_SignMismatch(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{
			name:  "negative value to unsigned target",
			input: -10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SafeConvert[uint](tt.input)
			if !errors.Is(err, ErrSignMismatch) {
				t.Errorf("SafeConvert() error = %v, want %v", err, ErrSignMismatch)
			}
		})
	}
}
