// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"context"

	commonCore "github.com/dedyf5/resik/core/common"
	reqOrganizationCore "github.com/dedyf5/resik/core/organization/request"
	resOrganizationCore "github.com/dedyf5/resik/core/organization/response"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	permisEntity "github.com/dedyf5/resik/entities/permission"
	"google.golang.org/grpc/codes"
)

func (h *OrganizationHandler) OrganizationPost(c context.Context, req *reqOrganizationCore.OrganizationPost) (*OrganizationUpsertRes, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("OrganizationPost")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	entity, err := req.ToEntity(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := h.organizationService.OrganizationInsert(ctx, entity); err != nil {
		return nil, err
	}

	return &OrganizationUpsertRes{
		Status: &commonCore.Status{
			Code: commonCore.StatusCodePlus(codes.OK),
			Message: term.SuccessfullyCreatedVal.Localize(
				ctx.Lang().Localizer,
				term.Organization.Localize(ctx.Lang().Localizer),
			),
		},
		Data: &resOrganizationCore.OrganizationUpsert{
			Id: entity.PublicID.String32(),
		},
	}, nil
}

func (h *OrganizationHandler) OrganizationPut(c context.Context, req *reqOrganizationCore.OrganizationPut) (*OrganizationUpsertRes, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("OrganizationPut")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	organizationID, _, err := ctx.GetOrganizationID(h.resolver, req.GetId(), permisEntity.PERMISSION_CODE_ORG_UPDATE.String())
	if err != nil {
		return nil, err
	}

	entity, err := req.ToEntity(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if _, err = h.organizationService.OrganizationUpdate(ctx, entity); err != nil {
		return nil, err
	}

	return &OrganizationUpsertRes{
		Status: &commonCore.Status{
			Code: commonCore.StatusCodePlus(codes.OK),
			Message: term.SuccessfullyUpdatedVal.Localize(
				ctx.Lang().Localizer,
				term.Organization.Localize(ctx.Lang().Localizer),
			),
		},
		Data: &resOrganizationCore.OrganizationUpsert{
			Id: req.GetId(),
		},
	}, nil
}
