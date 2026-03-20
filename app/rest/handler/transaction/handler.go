// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"net/http"

	echoFW "github.com/dedyf5/resik/app/rest/fw/echo"
	"github.com/dedyf5/resik/config"
	trxService "github.com/dedyf5/resik/core/transaction"
	reqTrxCore "github.com/dedyf5/resik/core/transaction/request"
	resTrxCore "github.com/dedyf5/resik/core/transaction/response"
	"github.com/dedyf5/resik/ctx"
	logCtx "github.com/dedyf5/resik/ctx/log"
	commonEntity "github.com/dedyf5/resik/entities/common"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v4"
)

// necessary to avoid unused package errors
// commonEntity package is used by swagger
var _ = commonEntity.Request{}

type Handler struct {
	config     config.Config
	log        *logCtx.Log
	fw         echoFW.IEcho
	trxService trxService.IService
}

func New(fw echoFW.IEcho, log *logCtx.Log, trxService trxService.IService, config config.Config) *Handler {
	return &Handler{
		config:     config,
		log:        log,
		fw:         fw,
		trxService: trxService,
	}
}

// @Summary Get Merchant Omzet
// @Description Get merchant omzet by merchant id
// @Tags transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       merchant_id path int true "Merchant ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       parameter query reqTrxCore.MerchantOmzetGet true "Query Param"
// @Success		200	{object}	resPkg.Response{data=[]resTrxCore.MerchantOmzet}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/transaction/merchant/{merchant_id}/omzet [get]
func (h *Handler) MerchantOmzetGet(echoCtx echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("MerchantOmzetGet")

	var payload reqTrxCore.MerchantOmzetGet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	param, err := payload.ToParam(ctx)
	if err != nil {
		return err
	}

	if _, err := ctx.UserClaims().MerchantIDIsAccessible(param.MerchantID); err != nil {
		return err
	}

	res, err := h.trxService.MerchantOmzetGet(param)
	if err != nil {
		return err
	}

	return &resPkg.Status{
		Code: http.StatusOK,
		Data: resTrxCore.MerchantOmzetFromEntity(res.Data),
		Meta: &resPkg.Meta{
			PageCurrent: param.Filter.Raw().PageOrDefault(),
			Limit:       param.Filter.Raw().LimitOrDefault(),
			Total:       res.Total,
		},
	}
}

// @Summary Get Outlet Omzet
// @Description Get outlet omzet by outlet id
// @Tags transaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       outlet_id path int true "Outlet ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       parameter query reqTrxCore.OutletOmzetGet true "Query Param"
// @Success		200	{object}	resPkg.Response{data=[]resTrxCore.OutletOmzet}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/transaction/outlet/{outlet_id}/omzet [get]
func (h *Handler) OutletOmzetGet(echoCtx echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}

	var payload reqTrxCore.OutletOmzetGet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	param, err := payload.ToParam(ctx)
	if err != nil {
		return err
	}

	if _, err := ctx.UserClaims().OutletIDIsAccessible(param.OutletID); err != nil {
		return err
	}

	res, err := h.trxService.OutletOmzetGet(param)
	if err != nil {
		return err
	}

	return &resPkg.Status{
		Code: http.StatusOK,
		Data: resTrxCore.OutletOmzetFromEntity(res.Data),
		Meta: &resPkg.Meta{
			PageCurrent: param.Filter.Raw().PageOrDefault(),
			Limit:       param.Filter.Raw().LimitOrDefault(),
			Total:       res.Total,
		},
	}
}
