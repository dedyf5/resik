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
	"github.com/dedyf5/resik/ctx/lang/term"
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

// TODO: Remove this dummy assignment once commonEntity is used explicitly elsewhere in this file.
// Currently kept to ensure Swagger can discover types from this package.
var _ = commonEntity.Request{}

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
// @Success		201	{object}	resPkg.ResponseSuccess{data=resMerchantCore.MerchantUpsert}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/merchant [post]
func (h *Handler) MerchantPost(echoCtx echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("MerchantPost")

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

	return resPkg.NewStatusSuccess(
		http.StatusCreated,
		term.SuccessfullyCreatedVal.Localize(
			ctx.Lang().Localizer,
			term.Merchant.Localize(ctx.Lang().Localizer),
		),
		&resMerchantCore.MerchantUpsert{
			Id: entity.ID,
		},
	)
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
// @Success		200	{object}	resPkg.ResponseSuccess{data=resMerchantCore.MerchantUpsert}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/merchant/{id} [put]
func (h *Handler) MerchantPut(echoCtx echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("MerchantPut")

	var body reqMerchantCore.MerchantPut

	if err := h.fw.StructValidator(echoCtx, &body); err != nil {
		return err
	}

	entity, err := body.ToEntity(ctx)
	if err != nil {
		return err
	}

	if _, err = ctx.UserClaims().MerchantIDIsAccessible(entity.ID); err != nil {
		return err
	}

	_, err = h.merchantService.MerchantUpdate(ctx, entity)
	if err != nil {
		return err
	}

	return resPkg.NewStatusSuccess(
		http.StatusOK,
		term.SuccessfullyUpdatedVal.Localize(
			ctx.Lang().Localizer,
			term.Merchant.Localize(ctx.Lang().Localizer),
		),
		&resMerchantCore.MerchantUpsert{
			Id: entity.ID,
		},
	)
}

// @Summary Get Merchant by ID
// @Description Get merchant by ID
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       id path int true "Merchant ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Success		200	{object}	resPkg.ResponseSuccess{data=resMerchantCore.MerchantDetail}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/merchant/{id} [get]
func (h *Handler) MerchantDetailGet(echoCtx echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("MerchantDetailGet")

	var param reqMerchantCore.MerchantDetailGet
	if err := h.fw.StructValidator(echoCtx, &param); err != nil {
		return err
	}

	if _, err = ctx.UserClaims().MerchantIDIsAccessible(param.GetId()); err != nil {
		return err
	}

	merchant, err := h.merchantService.MerchantGetByIDAndUserID(ctx, param.GetId(), ctx.UserClaims().UserID)
	if err != nil {
		return err
	}

	if merchant == nil {
		localizer := ctx.Lang().Localizer
		return resPkg.NewStatusMessage(
			http.StatusNotFound,
			term.NotFoundVal.Localize(localizer, term.Merchant.Localize(localizer)),
			nil,
		)
	}

	return resPkg.NewStatusData(
		http.StatusOK,
		resMerchantCore.MerchantDetailFromEntity(merchant),
	)
}

// @Summary Merchant List
// @Description Merchant list
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       parameter query reqMerchantCore.MerchantListGet true "Query Param"
// @Success		200	{object}	resPkg.ResponseSuccessWithMeta{data=[]resMerchantCore.MerchantList}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/merchant [get]
func (h *Handler) MerchantListGet(echoCtx echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("MerchantListGet")

	var payload reqMerchantCore.MerchantListGet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	param := payload.ToParam(ctx)

	res, err := h.merchantService.MerchantListGet(param)
	if err != nil {
		return err
	}

	return resPkg.NewStatusDataMeta(
		http.StatusOK,
		resMerchantCore.MerchantListFromEntity(res.Data),
		&resPkg.Meta{
			PageCurrent: param.Filter.Raw().PageOrDefault(),
			Limit:       param.Filter.Raw().LimitOrDefault(),
			Total:       res.Total,
		},
	)
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
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/merchant/{id} [delete]
func (h *Handler) MerchantDelete(echoCtx echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("MerchantDelete")

	var param reqMerchantCore.MerchantDelete
	if err := h.fw.StructValidator(echoCtx, &param); err != nil {
		return err
	}

	if _, err = ctx.UserClaims().MerchantIDIsAccessible(param.GetId()); err != nil {
		return err
	}

	_, err = h.merchantService.MerchantDelete(ctx, param.ToMerchant())
	if err != nil {
		return err
	}

	return resPkg.NewStatusCode(http.StatusNoContent)
}
