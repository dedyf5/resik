// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package handler

import (
	"net/http"

	echoFW "github.com/dedyf5/resik/app/http/fw/echo"
	trxReq "github.com/dedyf5/resik/app/http/handler/transaction/request"
	trxRes "github.com/dedyf5/resik/app/http/handler/transaction/response"
	"github.com/dedyf5/resik/config"
	trxService "github.com/dedyf5/resik/core/transaction"
	"github.com/dedyf5/resik/ctx"
	logCtx "github.com/dedyf5/resik/ctx/log"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	fw      echoFW.IEcho
	log     *logCtx.Log
	service trxService.IService
	config  config.Config
}

func New(fw echoFW.IEcho, log *logCtx.Log, service trxService.IService, config config.Config) *Handler {
	return &Handler{
		fw:      fw,
		log:     log,
		service: service,
		config:  config,
	}
}

// @Summary Get Merchant Omzet
// @Description Get merchant omzet by merchant id
// @Tags transaction
// @Accept json
// @Produce json
// @Param       id path int true "Merchant ID"
// @Param       parameter query request.MerchantOmzetGet true "Query Param"
// @Success		200	{object}	resPkg.Response{data=[]trxRes.MerchantOmzet}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/transaction/merchant/{id}/omzet [get]
func (h *Handler) MerchantOmzetGet(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	ctx.App.Logger().Debug("GetMerchantOmzet")

	var payload trxReq.MerchantOmzetGet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	param := payload.ToParam(ctx)

	res, err := h.service.MerchantOmzetGet(param)
	if err != nil {
		return err
	}

	return &resPkg.Status{
		Code: http.StatusOK,
		Data: trxRes.MerchantOmzetFromEntity(res.Data),
		Meta: &resPkg.Meta{
			PageCurrent: param.Filter.Page,
			Limit:       param.Filter.Limit,
			Total:       res.Total,
		},
	}
}

// @Summary Get Outlet Omzet
// @Description Get outlet omzet by outlet id
// @Tags transaction
// @Accept json
// @Produce json
// @Param       id path int true "Outlet ID"
// @Param       parameter query request.OutletOmzetGet true "Query Param"
// @Success		200	{object}	resPkg.Response{data=[]trxRes.OutletOmzet}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/transaction/outlet/{id}/omzet [get]
func (h *Handler) OutletOmzetGet(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}

	var payload trxReq.OutletOmzetGet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	param := payload.ToParam(ctx)

	res, err := h.service.OutletOmzetGet(param)
	if err != nil {
		return err
	}

	return &resPkg.Status{
		Code: http.StatusOK,
		Data: trxRes.OutletOmzetFromEntity(res.Data),
		Meta: &resPkg.Meta{
			PageCurrent: param.Filter.Page,
			Limit:       param.Filter.Limit,
			Total:       res.Total,
		},
	}
}
