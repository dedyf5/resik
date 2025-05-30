// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package echo

import (
	"fmt"
	"net/http"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	resPkg "github.com/dedyf5/resik/pkg/response"
	httpUtil "github.com/dedyf5/resik/utils/http"
	validatorUtil "github.com/dedyf5/resik/utils/validator"
	"github.com/labstack/echo/v4"
)

const DocPrefix string = "/docs/swagger"

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
	case *resPkg.Status:
		ctx.JSON(res.Code, httpUtil.LoggerFromStatus(res))
		return
	case *echo.HTTPError:
		ctx.JSON(res.Code, httpUtil.LoggerFromStatus(&resPkg.Status{
			Code:    res.Code,
			Message: fmt.Sprintf("%s", res.Message),
		}))
		return
	}

	ctx.JSON(http.StatusInternalServerError, httpUtil.LoggerErrorAuto(&resPkg.Status{
		Code:       http.StatusInternalServerError,
		CauseError: err,
	}))
}
