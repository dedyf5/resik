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
	statusEntity "github.com/dedyf5/resik/entities/status"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	fw      echoFW.IEcho
	service service.IService
	config  config.Config
}

func New(fw echoFW.IEcho, service service.IService, config config.Config) *Handler {
	return &Handler{
		fw:      fw,
		service: service,
		config:  config,
	}
}

func (h *Handler) GetMerchantOmzet(ctx echo.Context) error {
	var payload request.GeMerchantOmzet

	if err := h.fw.StructValidator(ctx, &payload); err != nil {
		return err
	}

	return &statusEntity.HTTP{
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
