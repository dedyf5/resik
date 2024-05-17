// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package general

import (
	"errors"
	"net/http"

	echoFW "github.com/dedyf5/resik/app/rest/fw/echo"
	"github.com/dedyf5/resik/config"
	resAppCore "github.com/dedyf5/resik/core/app/response"
	"github.com/dedyf5/resik/ctx"
	logCtx "github.com/dedyf5/resik/ctx/log"
	commonEntity "github.com/dedyf5/resik/entities/common"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	config config.Config
	log    *logCtx.Log
	fw     echoFW.IEcho
}

func New(config config.Config, log *logCtx.Log, fw echoFW.IEcho) *Handler {
	return &Handler{
		config: config,
		log:    log,
		fw:     fw,
	}
}

// @Summary Get Home
// @Description App info
// @Tags home
// @Accept json
// @Produce json
// @Param       parameter query commonEntity.Request true "Query Param"
// @Success		200	{object}	resPkg.Response{data=resAppCore.App}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/ [get]
func (h *Handler) Home(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}

	return &resPkg.Status{
		Code:       http.StatusInternalServerError,
		Message:    "error untuk menungso",
		CauseError: errors.New("error dari native"),
	}

	var payload commonEntity.Request

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	lang := ctx.Lang
	return &resPkg.Status{
		Code:    http.StatusOK,
		Message: lang.GetByTemplateData("home_message", commonEntity.Map{"app_name": h.config.App.Name, "code": h.config.App.Version}),
		Data:    resAppCore.AppMap(ctx, h.config),
	}
}
