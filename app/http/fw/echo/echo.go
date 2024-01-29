// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"errors"
	"net/http"

	"github.com/dedyf5/resik/ctx"
	httpApp "github.com/dedyf5/resik/ctx/app/http"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	httpUtil "github.com/dedyf5/resik/utils/http"
	logUtil "github.com/dedyf5/resik/utils/log"
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
	if err := e.validator.Struct(payload); err != nil {
		return err
	}
	return nil
}

func NewCtx(echoCtx echo.Context, log *logUtil.Log) (cxt *ctx.Ctx, err *httpApp.Status) {
	val := echoCtx.Get(langCtx.ContextKey.String())
	langRes, ok := val.(*langCtx.Lang)
	if !ok {
		return nil, &httpApp.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errors.New("failed to casting lang from context"),
		}
	}

	return ctx.NewHTTP(echoCtx.Request().Context(), log, langRes, echoCtx.Request().RequestURI), nil
}

func HTTPErrorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	switch status := err.(type) {
	case *httpApp.Status:
		if status.Code != http.StatusNoContent {
			ctx.JSON(status.Code, httpUtil.ResponseFromStatusHTTP(status))
		}
		return
	}

	ctx.JSON(http.StatusNotFound, httpUtil.ResponseErrorAuto(&httpApp.Status{
		Code: http.StatusNotFound,
	}))
}
