// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package handler

import (
	"net/http"

	echoFW "github.com/dedyf5/resik/app/http/fw/echo"
	"github.com/dedyf5/resik/app/http/handler/general/response"
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/ctx"
	logCtx "github.com/dedyf5/resik/ctx/log"
	commonEntity "github.com/dedyf5/resik/entities/common"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	fw     echoFW.IEcho
	log    *logCtx.Log
	config config.Config
}

func New(fw echoFW.IEcho, log *logCtx.Log, config config.Config) *Handler {
	return &Handler{
		fw:     fw,
		log:    log,
		config: config,
	}
}

// @Summary Get Home
// @Description App info
// @Tags home
// @Accept json
// @Produce json
// @Param       parameter query commonEntity.Request true "Query Param"
// @Success		200	{object}	resPkg.Response{data=response.Home}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/ [get]
func (h *Handler) Home(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}

	var payload commonEntity.Request

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	lang := ctx.Lang
	return &resPkg.Status{
		Code:    http.StatusOK,
		Message: lang.GetByTemplateData("home_message", commonEntity.Map{"app_name": h.config.App.Name, "code": h.config.App.Version}),
		Data:    response.HomeMap(ctx, h.config),
	}
}
