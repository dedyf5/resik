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

	"github.com/google/uuid"
)

type UUIDV7 uuid.UUID

var (
	Nil UUIDV7
)

func NewUUIDV7() (UUIDV7, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return Nil, err
	}
	return UUIDV7(u), nil
}

func ParseUUIDV7(value string) (UUIDV7, error) {
	u, err := uuid.Parse(value)
	if err != nil {
		return Nil, err
	}
	return UUIDV7(u), nil
}

func (u UUIDV7) uuid() uuid.UUID {
	return uuid.UUID(u)
}

func (u UUIDV7) String() string {
	return u.uuid().String()
}

// String32 return string without dash
// Format: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx (32 chars)
func (u UUIDV7) String32() string {
	return strings.ReplaceAll(u.uuid().String(), "-", "")
}

// Equal compare string with uuid
func (u UUIDV7) Equal(s string) bool {
	return u.String32() == strings.ReplaceAll(s, "-", "")
}

// MarshalJSON convert uuid to json
//
// Format: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" (32 chars + 2 quotes)
// Example: "019deba3afa8794691acaed1ef0cef3b"
func (u UUIDV7) MarshalJSON() ([]byte, error) {
	res := make([]byte, 34)
	res[0] = '"'
	res[33] = '"'

	hex.Encode(res[1:33], u[:])

	return res, nil
}

// UnmarshalJSON convert json to uuid
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

// IsEmpty check if uuid is Nil
func (u UUIDV7) IsEmpty() bool {
	return u == Nil
}

// Value return driver value
func (u UUIDV7) Value() (driver.Value, error) {
	return u.uuid().String(), nil
}

// Scan convert value to uuid
func (u *UUIDV7) Scan(value any) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("uuid: invalid value")
	}

	parseUUID, err := uuid.Parse(string(bytes))
	if err != nil {
		return err
	}

	*u = UUIDV7(parseUUID)
	return nil
}
