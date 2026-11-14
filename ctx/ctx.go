// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package ctx provides a unified request-scoped context wrapper that aggregates
// standard Go context, localization settings, structured logging, authenticated user claims.
package ctx

import (
	"context"
	"fmt"
	"runtime"

	jwt "github.com/dedyf5/resik/ctx/jwt"
	lang "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	permisEntity "github.com/dedyf5/resik/entities/permission"
	"github.com/dedyf5/resik/internal/identity"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

// Ctx wraps the standard context.Context alongside request-scoped dependencies
// including localization, structured logging, and authenticated JWT claims.
type Ctx struct {
	Context    context.Context
	lang       *lang.Lang
	log        *logCtx.Log
	userClaims *jwt.AuthClaims
}

// NewCtx initializes a new Ctx instance from an incoming context.Context and logger.
// It captures the caller's file and line location for telemetry, extracts localization settings,
// and attaches authenticated user claims if available.
//
// Returns a *resPkg.Status error (HTTP 500) if localization initialization fails.
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

// Lang returns the request-scoped localization handler.
func (c *Ctx) Lang() *lang.Lang {
	return c.lang
}

// Log returns the request-scoped structured logger instance.
func (c *Ctx) Log() *logCtx.Log {
	return c.log
}

// UserClaims returns the authenticated JWT claims associated with the current request.
func (c *Ctx) UserClaims() *jwt.AuthClaims {
	return c.userClaims
}

// GetOrganizationID resolves the internal organization ID from its public UUIDv7 string
// and verifies whether the authenticated user holds the specified permission code for that organization.
func (c *Ctx) GetOrganizationID(resolver identity.IdentityResolver, organizationPublicID string, permissionCode permisEntity.PERMISSION_CODE) (organizationID uint64, publicID uuidPkg.UUIDV7, err *resPkg.Status) {
	return c.userClaims.GetOrganizationID(c.Context, resolver, c.lang, organizationPublicID, permissionCode)
}

// GetBranchID resolves the internal branch ID from its public UUIDv7 string
// and verifies whether the authenticated user holds the specified permission code for that branch.
func (c *Ctx) GetBranchID(resolver identity.IdentityResolver, branchPublicID string, permissionCode permisEntity.PERMISSION_CODE) (branchID uint64, publicID uuidPkg.UUIDV7, err *resPkg.Status) {
	return c.userClaims.GetBranchID(c.Context, resolver, c.lang, branchPublicID, permissionCode)
}
