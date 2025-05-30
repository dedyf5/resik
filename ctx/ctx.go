// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package ctx

import (
	"context"
	"errors"
	"net/http"

	"github.com/dedyf5/resik/ctx/app"
	httpApp "github.com/dedyf5/resik/ctx/app/http"
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
	App        app.IApp
	Lang       *lang.Lang
	Context    context.Context
	UserClaims *jwt.AuthClaims
}

// return *Ctx HTTP. if create failed return *status.Status error
//
// error status code: 500
func NewHTTP(ctx context.Context, log *logCtx.Log, uri string) (*Ctx, *resPkg.Status) {
	langRes, err := lang.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &Ctx{
		App:        httpApp.NewHTTP(log.FromContext(ctx), uri),
		Context:    ctx,
		Lang:       langRes,
		UserClaims: jwt.AuthClaimsFromContext(ctx),
	}, nil
}

// return *Ctx HTTP. if create failed return *status.Status error
//
// error status code: 500
func NewHTTPFromGRPC(c context.Context, log *logCtx.Log) (*Ctx, *resPkg.Status) {
	if fullMethod, ok := c.Value(KeyFullMethod).(string); ok {
		return NewHTTP(c, log, fullMethod)
	}
	return nil, &resPkg.Status{
		Code:       http.StatusInternalServerError,
		CauseError: errors.New("failed to casting " + string(KeyFullMethod) + " from context"),
	}
}

func (c *Ctx) IsError() bool {
	return c.App.Status().IsError()
}

func (c *Ctx) Error() string {
	return c.App.Status().Error()
}

func (c *Ctx) MessageOrDefault() string {
	return c.App.Status().MessageOrDefault()
}
