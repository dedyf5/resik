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

// GetOrganizationID gets the organization ID by organization public ID, and check if user has access to it for permission code
func (c *Ctx) GetOrganizationID(resolver identity.IdentityResolver, organizationPublicID string, permissionCode string) (organizationID uint64, err *resPkg.Status) {
	return c.userClaims.GetOrganizationID(c.Context, resolver, c.lang, organizationPublicID, permissionCode)
}

// GetBranchID gets the branch ID by branch public ID, and check if user has access to it for permission code
func (c *Ctx) GetBranchID(resolver identity.IdentityResolver, branchPublicID string, permissionCode string) (branchID uint64, err *resPkg.Status) {
	return c.userClaims.GetBranchID(c.Context, resolver, c.lang, branchPublicID, permissionCode)
}
