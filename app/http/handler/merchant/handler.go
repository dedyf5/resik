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
// @Success		200	{object}	resPkg.Response{data=response.MerchantUpsert}
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
		Code: http.StatusOK,
		Message: ctx.Lang.GetByTemplateData("successfully_created_val", commonEntity.Map{
			"val": "merchant",
		}),
		Data: &response.MerchantUpsert{
			ID: entity.ID,
		},
	}
}
