// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"github.com/dedyf5/resik/ctx"
	userParam "github.com/dedyf5/resik/entities/user/param"
)

func (l *LoginPost) ToParam(c *ctx.Ctx) userParam.Auth {
	return userParam.Auth{
		Ctx:      c,
		Username: l.Username,
		Password: l.Password,
	}
}
