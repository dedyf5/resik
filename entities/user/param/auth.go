// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package param

import (
	"github.com/dedyf5/resik/ctx"
)

type Auth struct {
	Ctx      *ctx.Ctx
	Username string
	Password string
}
