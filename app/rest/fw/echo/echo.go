// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package echo provides Echo v5 framework helpers, request binding/validation,
// custom HTTP error handling, and REST middlewares for the application.
package echo

import (
	"errors"
	"log"
	"net/http"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	resPkg "github.com/dedyf5/resik/pkg/response"
	httpUtil "github.com/dedyf5/resik/utils/http"
	validatorUtil "github.com/dedyf5/resik/utils/validator"
	"github.com/labstack/echo/v5"
)

// DocPrefix defines the URL path prefix used to identify Swagger documentation routes.
const DocPrefix string = "/docs/swagger"

// IEcho defines the contract for Echo-specific request payload validation.
type IEcho interface {
	StructValidator(ctx *echo.Context, payload any) error
}

// Echo wraps a structured validator implementation for use with Echo request contexts.
type Echo struct {
	validator validatorUtil.IValidate
}

// New constructs a new Echo helper instance configured with the given struct validator.
func New(validator validatorUtil.IValidate) *Echo {
	return &Echo{
		validator: validator,
	}
}

// StructValidator binds the incoming request payload and validates struct fields,
// using the language resolved from the request context.
func (e *Echo) StructValidator(ctx *echo.Context, payload any) error {
	if err := ctx.Bind(payload); err != nil {
		return err
	}
	lang, err := langCtx.FromContext(ctx.Request().Context())
	if err != nil {
		return err
	}
	if err := e.validator.Struct(payload, lang); err != nil {
		return err
	}
	return nil
}

// HTTPErrorHandler formats errors returned by Echo handlers or middlewares into
// standardized JSON responses and logs committed-response or unexpected errors.
func HTTPErrorHandler(ctx *echo.Context, err error) {
	resp, errUnwrap := echo.UnwrapResponse(ctx.Response())
	if errUnwrap != nil {
		log.Printf("[error echo.UnwrapResponse] %s %s: %s\n", ctx.Request().Method, ctx.Request().URL.Path, errUnwrap.Error())
	}

	if resp != nil && resp.Committed {
		return
	}

	if res, ok := errors.AsType[*resPkg.Status](err); ok {
		_ = ctx.JSON(res.Code, httpUtil.LoggerFromStatus(res))
		return
	}

	var sc echo.HTTPStatusCoder
	if errors.As(err, &sc) {
		statusRes := resPkg.NewStatusCode(sc.StatusCode())
		_ = ctx.JSON(sc.StatusCode(), httpUtil.LoggerFromStatus(statusRes))
		return
	}

	statusRes := resPkg.NewStatusError(http.StatusInternalServerError, err)

	_ = ctx.JSON(http.StatusInternalServerError, httpUtil.LoggerErrorAuto(statusRes))
}
