// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"net/http"

	echoFW "github.com/dedyf5/resik/app/rest/fw/echo"
	merchantCore "github.com/dedyf5/resik/core/merchant"
	reqMerchantCore "github.com/dedyf5/resik/core/merchant/request"
	resMerchantCore "github.com/dedyf5/resik/core/merchant/response"
	"github.com/dedyf5/resik/ctx"
	logCtx "github.com/dedyf5/resik/ctx/log"
	commonEntity "github.com/dedyf5/resik/entities/common"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	log             *logCtx.Log
	fw              echoFW.IEcho
	merchantService merchantCore.IService
}

func New(log *logCtx.Log, fw echoFW.IEcho, merchantService merchantCore.IService) *Handler {
	return &Handler{
		log:             log,
		fw:              fw,
		merchantService: merchantService,
	}
}

// @Summary Create Merchant
// @Description Create new merchant
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       payload body reqMerchantCore.MerchantPost true "Payload"
// @Success		201	{object}	resPkg.Response{data=resMerchantCore.MerchantUpsert}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/merchant [post]
func (h *Handler) MerchantPost(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	ctx.App.Logger().Debug("MerchantPost")

	var payload reqMerchantCore.MerchantPost

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	entity, err := payload.ToEntity(ctx)
	if err != nil {
		return err
	}

	_, err = h.merchantService.MerchantInsert(ctx, entity)
	if err != nil {
		return err
	}

	return &resPkg.Status{
		Code: http.StatusCreated,
		Message: ctx.Lang.GetByTemplateData("successfully_created_val", commonEntity.Map{
			"val": "merchant",
		}),
		Data: &resMerchantCore.MerchantUpsert{
			ID: entity.ID,
		},
	}
}

// @Summary Update Merchant
// @Description Update merchant
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       id path int true "Merchant ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       payload body reqMerchantCore.MerchantPut true "Payload"
// @Success		200	{object}	resPkg.Response{data=resMerchantCore.MerchantUpsert}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     401 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/merchant/{id} [put]
func (h *Handler) MerchantPut(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	ctx.App.Logger().Debug("MerchantPut")

	var body reqMerchantCore.MerchantPut

	if err := h.fw.StructValidator(echoCtx, &body); err != nil {
		return err
	}

	entity, err := body.ToEntity(ctx)
	if err != nil {
		return err
	}

	if _, err = ctx.UserClaims.MerchantIDIsAccessible(entity.ID); err != nil {
		return err
	}

	_, err = h.merchantService.MerchantUpdate(ctx, entity)
	if err != nil {
		return err
	}

	return &resPkg.Status{
		Code: http.StatusOK,
		Message: ctx.Lang.GetByTemplateData("successfully_updated_val", commonEntity.Map{
			"val": "merchant",
		}),
		Data: &resMerchantCore.MerchantUpsert{
			ID: entity.ID,
		},
	}
}

// @Summary Merchant List
// @Description Merchant list
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       parameter query reqMerchantCore.MerchantListGet true "Query Param"
// @Success		200	{object}	resPkg.Response{data=resMerchantCore.MerchantUpsert}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     404 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/merchant [get]
func (h *Handler) MerchantListGet(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	ctx.App.Logger().Debug("MerchantListGet")

	var payload reqMerchantCore.MerchantListGet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	param := payload.ToParam(ctx)

	res, err := h.merchantService.MerchantListGet(param)
	if err != nil {
		return err
	}

	code := http.StatusOK
	if len(res.Data) == 0 {
		code = http.StatusNotFound
	}

	return &resPkg.Status{
		Code: code,
		Data: resMerchantCore.MerchantListFromEntity(res.Data),
		Meta: &resPkg.Meta{
			PageCurrent: param.Filter.Page,
			Limit:       param.Filter.Limit,
			Total:       res.Total,
		},
	}
}

// @Summary Update Merchant
// @Description Update merchant
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       id path int true "Merchant ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       parameter query reqMerchantCore.MerchantDelete true "Query Param"
// @Success		204	{object}	nil
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     401 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/merchant/{id} [delete]
func (h *Handler) MerchantDelete(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	ctx.App.Logger().Debug("MerchantDelete")

	var param reqMerchantCore.MerchantDelete
	if err := h.fw.StructValidator(echoCtx, &param); err != nil {
		return err
	}

	if _, err = ctx.UserClaims.MerchantIDIsAccessible(param.ID); err != nil {
		return err
	}

	_, err = h.merchantService.MerchantDelete(ctx, param.ToMerchant())
	if err != nil {
		return err
	}

	return &resPkg.Status{
		Code: http.StatusNoContent,
	}
}
