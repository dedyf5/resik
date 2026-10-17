// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package uuid

import (
	"testing"
	"uuid"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name            string
		inputUUID       uuid.UUID
		expectedVersion int
	}{
		{
			name: "UUID Version 4",
			inputUUID: uuid.UUID{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x40, // Byte 6: high nibble 4 represents version 4
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			expectedVersion: 4,
		},
		{
			name: "UUID Version 7",
			inputUUID: uuid.UUID{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x71, // Byte 6: high nibble 7 with non-zero low nibble
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			expectedVersion: 7,
		},
		{
			name: "UUID Version 1",
			inputUUID: uuid.UUID{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x1f, // Byte 6: high nibble 1 with max low nibble
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			expectedVersion: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Version(tt.inputUUID)
			if got != tt.expectedVersion {
				t.Errorf("Version() = %d; want %d", got, tt.expectedVersion)
			}
		})
	}
}
