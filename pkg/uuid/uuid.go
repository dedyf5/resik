// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package uuid

import "uuid"

// Version returns the version of a UUID.
func Version(u uuid.UUID) int {
	return int(u[6] >> 4)
}
