// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package ctx

import (
	"context"
	"fmt"
	"runtime"

	jwt "github.com/dedyf5/resik/ctx/jwt"
	lang "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/internal/identity"
	resPkg "github.com/dedyf5/resik/pkg/response"
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
	_, file, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s:%d", file, line)

	if holder, ok := c.Value(logCtx.KeyCallerHolderContext).(*logCtx.CallerHolder); ok {
		holder.Caller = caller
	}

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

// GetMerchantID gets the merchant ID by merchant public ID, and check if user has access to it
func (c *Ctx) GetMerchantID(resolver identity.IdentityResolver, merchantPublicID string) (merchantID uint64, err *resPkg.Status) {
	return c.userClaims.GetMerchantID(c.Context, resolver, c.lang, merchantPublicID)
}

// GetOutletID gets the outlet ID by outlet public ID, and check if user has access to it
func (c *Ctx) GetOutletID(resolver identity.IdentityResolver, outletPublicID string) (outletID uint64, err *resPkg.Status) {
	return c.userClaims.GetOutletID(c.Context, resolver, c.lang, outletPublicID)
}
