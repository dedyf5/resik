// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"net/http"

	echoFW "github.com/dedyf5/resik/app/http/fw/echo"
	"github.com/dedyf5/resik/app/http/handler/merchant/request"
	"github.com/dedyf5/resik/app/http/handler/merchant/response"
	"github.com/dedyf5/resik/config"
	merchantCore "github.com/dedyf5/resik/core/merchant"
	"github.com/dedyf5/resik/ctx"
	logCtx "github.com/dedyf5/resik/ctx/log"
	commonEntity "github.com/dedyf5/resik/entities/common"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	fw              echoFW.IEcho
	log             *logCtx.Log
	merchantService merchantCore.IService
	config          config.Config
}

func New(fw echoFW.IEcho, log *logCtx.Log, merchantService merchantCore.IService, config config.Config) *Handler {
	return &Handler{
		fw:              fw,
		log:             log,
		merchantService: merchantService,
		config:          config,
	}
}

// @Summary Create Merchant
// @Description Create new merchant
// @Tags merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       payload body request.MerchantPost true "Payload"
// @Success		201	{object}	resPkg.Response{data=response.MerchantUpsert}
// @Failure     400 {object}	resPkg.Response{data=nil}
// @Failure     500 {object}	resPkg.Response{data=nil}
// @Router		/merchant [post]
func (h *Handler) MerchantPost(echoCtx echo.Context) error {
	ctx, err := ctx.NewHTTP(echoCtx.Request().Context(), h.log, echoCtx.Request().RequestURI)
	if err != nil {
		return err
	}
	ctx.App.Logger().Debug("MerchantPost")

	var payload request.MerchantPost

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	var query commonEntity.Request
	if err := h.fw.StructValidator(echoCtx, &query); err != nil {
		return err
	}

	entity, err := request.MerchantPostToEntity(ctx, &payload)
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
		Data: &response.MerchantUpsert{
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
// @Param       payload body request.MerchantPut true "Payload"
// @Success		200	{object}	resPkg.Response{data=response.MerchantUpsert}
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

	var body request.MerchantPut

	if err := h.fw.StructValidator(echoCtx, &body); err != nil {
		return err
	}

	var param commonEntity.Request
	if err := h.fw.StructValidator(echoCtx, &param); err != nil {
		return err
	}

	entity, err := request.MerchantPutToEntity(ctx, &body)
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
		Data: &response.MerchantUpsert{
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
// @Param       parameter query request.MerchantListGet true "Query Param"
// @Success		200	{object}	resPkg.Response{data=response.MerchantUpsert}
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

	var payload request.MerchantListGet

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
		Data: response.MerchantListFromEntity(res.Data),
		Meta: &resPkg.Meta{
			PageCurrent: param.Filter.Page,
			Limit:       param.Filter.Limit,
			Total:       res.Total,
		},
	}
}
