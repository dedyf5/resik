// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package param defines input parameter structs for user authentication and related operations.
package param

import (
	"github.com/dedyf5/resik/ctx"
)

// Auth hold execution context and credentials for authenticating a user.
type Auth struct {
	Ctx      *ctx.Ctx
	Username string
	Password string
}
