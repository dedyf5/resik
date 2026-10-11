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
	permisEntity "github.com/dedyf5/resik/entities/permission"
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

// @Summary Get Organization by ID
// @Description Get organization by ID
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param       id path string true "Organization ID"
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

	organizationID, _, err := ctx.GetOrganizationID(h.resolver, param.GetId(), permisEntity.PERMISSION_CODE_ORG_READ.String())
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
