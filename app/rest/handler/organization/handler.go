// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"net/http"

	echoFW "github.com/dedyf5/resik/app/rest/fw/echo"
	organizationCore "github.com/dedyf5/resik/core/organization"
	reqOrganizationCore "github.com/dedyf5/resik/core/organization/request"
	resOrganizationCore "github.com/dedyf5/resik/core/organization/response"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	logCtx "github.com/dedyf5/resik/ctx/log"
	commonEntity "github.com/dedyf5/resik/entities/common"
	"github.com/dedyf5/resik/internal/identity"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	log                 *logCtx.Log
	fw                  echoFW.IEcho
	resolver            identity.IdentityResolver
	organizationService organizationCore.IService
}

// TODO: Remove this dummy assignment once commonEntity is used explicitly elsewhere in this file.
// Currently kept to ensure Swagger can discover types from this package.
var _ = commonEntity.Request{}

func New(log *logCtx.Log, fw echoFW.IEcho, resolver identity.IdentityResolver, organizationService organizationCore.IService) *Handler {
	return &Handler{
		log:                 log,
		fw:                  fw,
		resolver:            resolver,
		organizationService: organizationService,
	}
}

// @Summary Create Organization
// @Description Create new organization
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       payload body reqOrganizationCore.OrganizationPost true "Payload"
// @Success		201	{object}	resPkg.ResponseSuccess{data=resOrganizationCore.OrganizationUpsert}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/organizations [post]
func (h *Handler) OrganizationPost(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("OrganizationPost")

	var payload reqOrganizationCore.OrganizationPost

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	entity, err := payload.ToEntity(ctx)
	if err != nil {
		return err
	}

	_, err = h.organizationService.OrganizationInsert(ctx, entity)
	if err != nil {
		return err
	}

	return resPkg.NewStatusSuccess(
		http.StatusCreated,
		term.SuccessfullyCreatedVal.Localize(
			ctx.Lang().Localizer,
			term.Organization.Localize(ctx.Lang().Localizer),
		),
		&resOrganizationCore.OrganizationUpsert{
			Id: entity.PublicID.String32(),
		},
	)
}

// @Summary Update Organization
// @Description Update organization
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       id path int true "Organization ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       payload body reqOrganizationCore.OrganizationPut true "Payload"
// @Success		200	{object}	resPkg.ResponseSuccess{data=resOrganizationCore.OrganizationUpsert}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/organizations/{id} [put]
func (h *Handler) OrganizationPut(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("OrganizationPut")

	var body reqOrganizationCore.OrganizationPut

	if err := h.fw.StructValidator(echoCtx, &body); err != nil {
		return err
	}

	organizationID, err := ctx.GetOrganizationID(h.resolver, body.GetId(), "organization:update")
	if err != nil {
		return err
	}

	entity, err := body.ToEntity(ctx, organizationID)
	if err != nil {
		return err
	}

	_, err = h.organizationService.OrganizationUpdate(ctx, entity)
	if err != nil {
		return err
	}

	return resPkg.NewStatusSuccess(
		http.StatusOK,
		term.SuccessfullyUpdatedVal.Localize(
			ctx.Lang().Localizer,
			term.Organization.Localize(ctx.Lang().Localizer),
		),
		&resOrganizationCore.OrganizationUpsert{
			Id: body.GetId(),
		},
	)
}

// @Summary Get Organization by ID
// @Description Get organization by ID
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       id path int true "Organization ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Success		200	{object}	resPkg.ResponseSuccess{data=resOrganizationCore.OrganizationDetail}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/organizations/{id} [get]
func (h *Handler) OrganizationDetailGet(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("OrganizationDetailGet")

	var param reqOrganizationCore.OrganizationDetailGet
	if err := h.fw.StructValidator(echoCtx, &param); err != nil {
		return err
	}

	organizationID, err := ctx.GetOrganizationID(h.resolver, param.GetId(), "organization:read")
	if err != nil {
		return err
	}

	organization, err := h.organizationService.OrganizationGetByID(ctx, organizationID)
	if err != nil {
		return err
	}

	if organization == nil {
		localizer := ctx.Lang().Localizer
		return resPkg.NewStatusMessage(
			http.StatusNotFound,
			term.NotFoundVal.Localize(localizer, term.Organization.Localize(localizer)),
			nil,
		)
	}

	return resPkg.NewStatusData(
		http.StatusOK,
		resOrganizationCore.OrganizationDetailFromDTO(organization),
	)
}

// @Summary Organization List
// @Description Organization list
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       parameter query reqOrganizationCore.OrganizationListGet true "Query Param"
// @Success		200	{object}	resPkg.ResponseSuccessWithMeta{data=[]resOrganizationCore.OrganizationList}
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/organizations [get]
func (h *Handler) OrganizationListGet(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("OrganizationListGet")

	var payload reqOrganizationCore.OrganizationListGet

	if err := h.fw.StructValidator(echoCtx, &payload); err != nil {
		return err
	}

	param, err := payload.ToParam(ctx, h.resolver)
	if err != nil {
		return err
	}

	res, err := h.organizationService.OrganizationsGet(param)
	if err != nil {
		return err
	}

	return resPkg.NewStatusDataMeta(
		http.StatusOK,
		resOrganizationCore.OrganizationListFromDTO(&res.Data),
		&resPkg.Meta{
			PageCurrent: param.Filter.Raw().PageOrDefault(),
			Limit:       param.Filter.Raw().LimitOrDefault(),
			Total:       res.Total,
		},
	)
}

// @Summary Delete Organization
// @Description Delete organization
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       id path int true "Organization ID"
// @Param       parameter query commonEntity.Request true "Query Param"
// @Param       parameter query reqOrganizationCore.OrganizationDelete true "Query Param"
// @Success		204	{object}	nil
// @Failure     400 {object}	resPkg.ResponseBadRequest
// @Failure     401 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     429 {object}	resPkg.ResponseErrorWithoutDetails
// @Failure     500 {object}	resPkg.ResponseErrorWithoutDetails
// @Router		/organizations/{id} [delete]
func (h *Handler) OrganizationDelete(echoCtx *echo.Context) error {
	ctx, err := ctx.NewCtx(echoCtx.Request().Context(), h.log)
	if err != nil {
		return err
	}
	h.log.Debug("OrganizationDelete")

	var param reqOrganizationCore.OrganizationDelete
	if err := h.fw.StructValidator(echoCtx, &param); err != nil {
		return err
	}

	organizationID, err := ctx.GetOrganizationID(h.resolver, param.GetId(), "organization:delete")
	if err != nil {
		return err
	}

	_, err = h.organizationService.OrganizationDelete(ctx, param.ToOrganization(organizationID))
	if err != nil {
		return err
	}

	return resPkg.NewStatusCode(http.StatusNoContent)
}
