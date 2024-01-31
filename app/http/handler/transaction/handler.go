// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package handler

import (
	"errors"
	"net/http"

	echoFW "github.com/dedyf5/resik/app/http/fw/echo"
	"github.com/dedyf5/resik/app/http/handler/transaction/request"
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/core/transaction/service"
	"github.com/dedyf5/resik/ctx"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/ctx/status"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	fw      echoFW.IEcho
	log     *logCtx.Log
	service service.IService
	config  config.Config
}

func New(fw echoFW.IEcho, log *logCtx.Log, service service.IService, config config.Config) *Handler {
	return &Handler{
		fw:      fw,
		log:     log,
		service: service,
		config:  config,
	}
}

func (h *Handler) GetMerchantOmzet(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	ctx.App.Logger().Debug("GetMerchantOmzet")

	var payload request.GeMerchantOmzet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	return &status.Status{
		Code: http.StatusOK,
		Data: map[string]interface{}{
			"page_or_default":  payload.PageOrDefault(),
			"limit_or_default": payload.LimitOrDefault(),
			"req":              payload,
			"lang_def":         ctx.Lang.LangDefault.String(),
			"lang_req":         ctx.Lang.LangReq.String(),
		},
	}
}

func (h *Handler) GetOutletOmzet(ctx echo.Context) error {
	return errors.New("MASUK GetOutletOmzet")
}

func (h *Handler) Login(ctx echo.Context) error {
	return errors.New("MASUK Login")
}
