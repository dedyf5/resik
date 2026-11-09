// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package uuid

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"uuid"
)

// UUIDV7 represents a UUID version 7 type with custom database and JSON serialization support.
type UUIDV7 uuid.UUID

var (
	// Nil represents a zero-value UUIDV7 (all bytes set to zero).
	Nil UUIDV7

	// ErrInvalidValue is returned when a scanned value type is invalid or unsupported.
	ErrInvalidValue = errors.New("uuid: invalid value")

	// ErrInvalidVersion is returned when a parsed UUID does not match version 7.
	ErrInvalidVersion = errors.New("uuid: invalid version")
)

// NewUUIDV7 generates and returns a new UUIDv7 instance.
func NewUUIDV7() UUIDV7 {
	return UUIDV7(uuid.NewV7())
}

// ParseUUIDV7 parses a string representation into a UUIDV7, returning an error if parsing fails or version is not 7.
func ParseUUIDV7(value string) (UUIDV7, error) {
	u, err := uuid.Parse(value)
	if err != nil {
		return Nil, err
	}
	if Version(u) != 7 {
		return Nil, ErrInvalidVersion
	}
	return UUIDV7(u), nil
}

// uuid converts UUIDV7 to the underlying standard uuid.UUID type.
func (u UUIDV7) uuid() uuid.UUID {
	return uuid.UUID(u)
}

// String returns the standard string representation of the UUIDV7 with hyphens.
func (u UUIDV7) String() string {
	return u.uuid().String()
}

// String32 returns the string representation of the UUIDV7 without hyphens.
// Format: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx (32 chars)
func (u UUIDV7) String32() string {
	return strings.ReplaceAll(u.uuid().String(), "-", "")
}

// Equal compares the UUIDV7 with another string representation (with or without hyphens).
func (u UUIDV7) Equal(s string) bool {
	return u.String32() == strings.ReplaceAll(s, "-", "")
}

// MarshalJSON implements the json.Marshaler interface to encode UUIDV7 as a 32-character hex string enclosed in quotes.
// Format: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" (32 chars + 2 quotes)
// Example: "019deba3afa8794691acaed1ef0cef3b"
func (u UUIDV7) MarshalJSON() ([]byte, error) {
	res := make([]byte, 34)
	res[0] = '"'
	res[33] = '"'

	hex.Encode(res[1:33], u[:])

	return res, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface to decode a JSON string into a UUIDV7.
func (u *UUIDV7) UnmarshalJSON(data []byte) error {
	if data == nil {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parseUUID, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*u = UUIDV7(parseUUID)
	return nil
}

// IsEmpty checks whether the UUIDV7 is equal to Nil.
func (u UUIDV7) IsEmpty() bool {
	return u == Nil
}

// Value implements the driver.Valuer interface for database serialization.
func (u UUIDV7) Value() (driver.Value, error) {
	return u.uuid().String(), nil
}

// Scan implements the sql.Scanner interface for value deserialization.
func (u *UUIDV7) Scan(value any) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return ErrInvalidValue
	}

	parseUUID, err := uuid.Parse(string(bytes))
	if err != nil {
		return err
	}

	*u = UUIDV7(parseUUID)
	return nil
}
