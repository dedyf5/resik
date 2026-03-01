// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package ctx

import (
	"context"

	jwt "github.com/dedyf5/resik/ctx/jwt"
	lang "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

type Key string

const (
	KeyFullMethod Key = "FullMethod"
)

type Ctx struct {
	Context    context.Context
	lang       *lang.Lang
	log        *logCtx.Log
	userClaims *jwt.AuthClaims
}

// return *Ctx HTTP. if create failed return *status.Status error
//
// error status code: 500
func NewCtx(c context.Context, log *logCtx.Log) (*Ctx, *resPkg.Status) {
	langRes, err := lang.FromContext(c)
	if err != nil {
		return nil, err
	}
	return &Ctx{
		Context:    c,
		lang:       langRes,
		log:        log,
		userClaims: jwt.AuthClaimsFromContext(c),
	}, nil
}

func (c *Ctx) Lang() *lang.Lang {
	return c.lang
}

func (c *Ctx) Log() *logCtx.Log {
	return c.log
}

func (c *Ctx) UserClaims() *jwt.AuthClaims {
	return c.userClaims
}
