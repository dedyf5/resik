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
	"github.com/dedyf5/resik/internal/identity"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v5"
)

// necessary to avoid unused package errors
// commonEntity package is used by swagger
var _ = commonEntity.Request{}

type Handler struct {
	config     config.Config
	log        *logCtx.Log
	fw         echoFW.IEcho
	resolver   identity.IdentityResolver
	trxService trxService.IService
}

func New(config config.Config, log *logCtx.Log, fw echoFW.IEcho, resolver identity.IdentityResolver, trxService trxService.IService) *Handler {
	return &Handler{
		config:     config,
		log:        log,
		fw:         fw,
		resolver:   resolver,
		trxService: trxService,
	}
}

// @Summary Get Organization Omzet
// @Description Get organization omzet by organization id
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       organization_id path int true "Organization ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       parameter query reqTrxCore.OrganizationOmzetGet true "Query Param"
// @Success		200	{object}	resPkg.ResponseSuccessWithMeta{data=[]resTrxCore.OrganizationOmzet}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/transactions/organization/{organization_id}/omzet [get]
func (h *Handler) OrganizationOmzetGet(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("OrganizationOmzetGet")

	var payload reqTrxCore.OrganizationOmzetGet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	_, organizationPublicID, err := ctx.GetOrganizationID(h.resolver, payload.GetOrganizationId(), "transaction:read")
	if err != nil {
		return err
	}

	param, err := payload.ToParam(ctx, organizationPublicID)
	if err != nil {
		return err
	}

	res, err := h.trxService.OrganizationOmzetGet(param)
	if err != nil {
		return err
	}

	return resPkg.NewStatusDataMeta(
		http.StatusOK,
		resTrxCore.OrganizationOmzetFromEntity(res.Data),
		&resPkg.Meta{
			PageCurrent: param.Filter.Raw().PageOrDefault(),
			Limit:       param.Filter.Raw().LimitOrDefault(),
			Total:       res.Total,
		},
	)
}

// @Summary Get Branch Omzet
// @Description Get branch omzet by branch id
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       branch_id path int true "Branch ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       parameter query reqTrxCore.BranchOmzetGet true "Query Param"
// @Success		200	{object}	resPkg.ResponseSuccessWithMeta{data=[]resTrxCore.BranchOmzet}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/transactions/branch/{branch_id}/omzet [get]
func (h *Handler) BranchOmzetGet(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}

	var payload reqTrxCore.BranchOmzetGet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	_, branchPublicID, err := ctx.GetBranchID(h.resolver, payload.GetBranchId(), "transaction:read")
	if err != nil {
		return err
	}

	param, err := payload.ToParam(ctx, branchPublicID)
	if err != nil {
		return err
	}

	res, err := h.trxService.BranchOmzetGet(param)
	if err != nil {
		return err
	}

	return resPkg.NewStatusDataMeta(
		http.StatusOK,
		resTrxCore.BranchOmzetFromEntity(res.Data),
		&resPkg.Meta{
			PageCurrent: param.Filter.Raw().PageOrDefault(),
			Limit:       param.Filter.Raw().LimitOrDefault(),
			Total:       res.Total,
		},
	)
}
