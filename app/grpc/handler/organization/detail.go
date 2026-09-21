// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"context"
	"net/http"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqOrganizationCore "github.com/dedyf5/resik/core/organization/request"
	resOrganizationCore "github.com/dedyf5/resik/core/organization/response"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"google.golang.org/grpc/codes"
)

func (h *OrganizationHandler) OrganizationDetailGet(c context.Context, req *reqOrganizationCore.OrganizationDetailGet) (*OrganizationDetailGetRes, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("OrganizationDetailGet")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	organizationID, _, err := ctx.GetOrganizationID(h.resolver, req.GetId(), "organization:read")
	if err != nil {
		return nil, err
	}

	organization, err := h.organizationService.OrganizationGetByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization == nil {
		localizer := ctx.Lang().Localizer
		return nil, resPkg.NewStatusMessage(
			http.StatusNotFound,
			term.NotFoundVal.Localize(localizer, term.Organization.Localize(localizer)),
			nil,
		)
	}

	return &OrganizationDetailGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: resOrganizationCore.OrganizationDetailFromDTO(organization),
	}, nil
}
