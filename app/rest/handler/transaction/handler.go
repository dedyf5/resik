// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"net/http"

	echoFW "github.com/dedyf5/resik/app/rest/fw/echo"
	"github.com/dedyf5/resik/config"
	commonCore "github.com/dedyf5/resik/core/common"
	trxService "github.com/dedyf5/resik/core/transaction"
	dtoTrx "github.com/dedyf5/resik/core/transaction/dto"
	reqTrxCore "github.com/dedyf5/resik/core/transaction/request"
	resTrxCore "github.com/dedyf5/resik/core/transaction/response"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
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

// @Summary Create Transaction
// @Description Create new transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       payload body reqTrxCore.TransactionPost true "Payload"
// @Success		201	{object}	resPkg.ResponseSuccess{data=commonCore.Id}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/transactions [post]
func (h *Handler) TransactionPost(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("TransactionPost")

	var payload reqTrxCore.TransactionPost
	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	entity, err := payload.ToEntity(ctx)
	if err != nil {
		return err
	}

	_, err = h.trxService.TransactionInsert(ctx, entity)
	if err != nil {
		return err
	}

	return resPkg.NewStatusSuccess(
		http.StatusCreated,
		term.SuccessfullyCreatedVal.Localize(
			ctx.Lang().Localizer,
			term.Transaction.Localize(ctx.Lang().Localizer),
		),
		&commonCore.Id{Id: entity.PublicID.String32()},
	)
}

// @Summary Update Transaction
// @Description Update a transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       id path string true "Transaction ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       payload body reqTrxCore.TransactionPut true "Payload"
// @Success		200	{object}	resPkg.ResponseSuccess{data=commonCore.Id}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     403 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     422 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/transactions/{id} [put]
func (h *Handler) TransactionPut(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("TransactionPut")

	var payload reqTrxCore.TransactionPut
	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	entity, err := payload.ToEntity(ctx)
	if err != nil {
		return err
	}

	_, err = h.trxService.TransactionUpdate(ctx, entity)
	if err != nil {
		return err
	}

	return resPkg.NewStatusSuccess(
		http.StatusOK,
		term.SuccessfullyUpdatedVal.Localize(
			ctx.Lang().Localizer,
			term.Transaction.Localize(ctx.Lang().Localizer),
		),
		&commonCore.Id{Id: payload.GetId()},
	)
}

// @Summary Get Transaction by ID
// @Description Get transaction by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       id path string true "Transaction ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Success		200	{object}	resPkg.ResponseSuccess{data=resTrxCore.TransactionDetail}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     404 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/transactions/{id} [get]
func (h *Handler) TransactionDetailGet(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("TransactionDetailGet")

	var payload reqTrxCore.Id
	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	param, err := payload.ToParam(ctx, h.resolver, dtoTrx.PERMISSION_CODE_READ)
	if err != nil {
		return err
	}

	transaction, err := h.trxService.TransactionGetByPublicID(param)
	if err != nil {
		return err
	}

	if transaction == nil {
		localizer := ctx.Lang().Localizer
		return resPkg.NewStatusMessage(
			http.StatusNotFound,
			term.NotFoundVal.Localize(localizer, term.Transaction.Localize(localizer)),
			nil,
		)
	}

	return resPkg.NewStatusData(
		http.StatusOK,
		resTrxCore.TransactionDetailFromDTO(transaction),
	)
}

// @Summary Transactions List
// @Description Get transactions list
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       parameter query reqTrxCore.TransactionListGet true "Query Param"
// @Success		200	{object}	resPkg.ResponseSuccessWithMeta{data=[]resTrxCore.TransactionList}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/transactions [get]
func (h *Handler) TransactionListGet(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("TransactionListGet")

	var payload reqTrxCore.TransactionListGet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	param, err := payload.ToParam(ctx, h.resolver)
	if err != nil {
		return err
	}

	res, err := h.trxService.TransactionsGet(param)
	if err != nil {
		return err
	}

	return resPkg.NewStatusDataMeta(
		http.StatusOK,
		resTrxCore.TransactionsGetFromDTO(&res.Data),
		&resPkg.Meta{
			PageCurrent: param.Filter.Raw().PageOrDefault(),
			Limit:       param.Filter.Raw().LimitOrDefault(),
			Total:       res.Total,
		},
	)
}

// @Summary Delete Transaction
// @Description Delete transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       id path string true "Transaction ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Success		204	{object}	nil
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     404 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/transactions/{id} [delete]
func (h *Handler) TransactionDelete(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("TransactionDelete")

	var payload reqTrxCore.Id
	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	param, err := payload.ToParam(ctx, h.resolver, dtoTrx.PERMISSION_CODE_DELETE)
	if err != nil {
		return err
	}

	if ok, err := h.trxService.TransactionDelete(param); err != nil || !ok {
		return err
	}

	return resPkg.NewStatusCode(http.StatusNoContent)
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
