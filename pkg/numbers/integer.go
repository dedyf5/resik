// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package numbers

import (
	"errors"
	"fmt"
)

// Integer defines a type constraint for all signed and unsigned integer types.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Target Integer
type Source Integer

var (
	// ErrOverflow is returned when the value exceeds the capacity of the target type.
	ErrOverflow = errors.New("integer overflow: value exceeds target type capacity")

	// ErrSignMismatch is returned when converting a negative value to an unsigned type.
	ErrSignMismatch = errors.New("sign mismatch: cannot convert negative value to unsigned type")
)

// SafeConvert performs a bounds-checked conversion between any two integer types.
// It returns an error if the conversion would result in data loss or sign mismatch.
func SafeConvert[T Target, S Source](val S) (T, error) {
	result := T(val)

	if S(result) != val {
		return 0, fmt.Errorf("%w: value %v does not fit in %T", ErrOverflow, val, result)
	}

	if val < 0 && result > 0 {
		return 0, fmt.Errorf("%w: value %v to %T", ErrSignMismatch, val, result)
	}

	return result, nil
}
