// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package uuid provides UUID generation, parsing, and custom type wrappers
// with support for JSON marshaling, value scanners, and version validation.
package uuid

import "uuid"

// Version extracts and returns the version number from the given UUID.
func Version(u uuid.UUID) int {
	return int(u[6] >> 4)
}
