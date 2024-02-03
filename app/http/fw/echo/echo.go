// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"net/http"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/status"
	httpUtil "github.com/dedyf5/resik/utils/http"
	validatorUtil "github.com/dedyf5/resik/utils/validator"
	"github.com/labstack/echo/v4"
)

type IEcho interface {
	StructValidator(ctx echo.Context, payload interface{}) error
}

type Echo struct {
	validator validatorUtil.IValidate
}

func New(validator validatorUtil.IValidate) *Echo {
	return &Echo{
		validator: validator,
	}
}

func (e *Echo) StructValidator(ctx echo.Context, payload interface{}) error {
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

func HTTPErrorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	switch res := err.(type) {
	case *status.Status:
		if res.Code != http.StatusNoContent {
			ctx.JSON(res.Code, httpUtil.LoggerFromStatusHTTP(res))
		}
		return
	}

	ctx.JSON(http.StatusInternalServerError, httpUtil.LoggerErrorAuto(&status.Status{
		Code: http.StatusInternalServerError,
	}))
}
