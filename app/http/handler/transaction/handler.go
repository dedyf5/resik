// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package handler

import (
	"errors"
	"net/http"

	echoFW "github.com/dedyf5/resik/app/http/fw/echo"
	"github.com/dedyf5/resik/app/http/handler/transaction/request"
	"github.com/dedyf5/resik/app/http/handler/transaction/response"
	"github.com/dedyf5/resik/config"
	trxService "github.com/dedyf5/resik/core/transaction"
	"github.com/dedyf5/resik/ctx"
	logCtx "github.com/dedyf5/resik/ctx/log"
	resEntity "github.com/dedyf5/resik/entities/response"
	"github.com/dedyf5/resik/pkg/status"
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
// @Success		200	{object}	resEntity.Response{data=[]response.MerchantOmzet}
// @Failure     400 {object}	resEntity.Response{}
// @Failure     500 {object}	resEntity.Response{}
// @Router		/transaction/merchant/{id}/omzet [get]
func (h *Handler) MerchantOmzetGet(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	ctx.App.Logger().Debug("GetMerchantOmzet")

	var payload request.MerchantOmzetGet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	param := payload.ToParam(ctx)

	res, err := h.service.MerchantOmzetGet(param)
	if err != nil {
		return err
	}

	return &status.Status{
		Code: http.StatusOK,
		Data: response.MerchantOmzetFromEntity(res.Data),
		Meta: &status.Meta{
			PageCurrent: param.Filter.Page,
			Limit:       param.Filter.Limit,
			Total:       res.Total,
		},
		Format: resEntity.Response{},
	}
}

func (h *Handler) Create(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	// ctx.App.Logger().Debug("Create")

	var payload request.TransactionUpsert

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	langReq := ""
	if ctx.Lang.LangReq != nil {
		langReq = ctx.Lang.LangReq.String()
	}

	return &status.Status{
		Code: http.StatusOK,
		Data: map[string]interface{}{
			"req":      payload,
			"lang_def": ctx.Lang.LangDefault.String(),
			"lang_req": langReq,
		},
	}
}

func (h *Handler) OutletOmzetGet(ctx echo.Context) error {
	return errors.New("MASUK GetOutletOmzet")
}

func (h *Handler) Login(ctx echo.Context) error {
	return errors.New("MASUK Login")
}
