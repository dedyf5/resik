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
	"github.com/dedyf5/resik/ctx/status"
	logUtil "github.com/dedyf5/resik/utils/log"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	fw      echoFW.IEcho
	log     *logUtil.Log
	service service.IService
	config  config.Config
}

func New(fw echoFW.IEcho, log *logUtil.Log, service service.IService, config config.Config) *Handler {
	return &Handler{
		fw:      fw,
		log:     log,
		service: service,
		config:  config,
	}
}

func (h *Handler) GetMerchantOmzet(echoCtx echo.Context) error {
	ctx, err := echoFW.NewCtx(echoCtx, h.log)
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
		},
	}
}

func (h *Handler) GetOutletOmzet(ctx echo.Context) error {
	return errors.New("MASUK GetOutletOmzet")
}

func (h *Handler) Login(ctx echo.Context) error {
	return errors.New("MASUK Login")
}
