// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqOrganizationCore "github.com/dedyf5/resik/core/organization/request"
	"github.com/dedyf5/resik/core/organization/response"
	"github.com/dedyf5/resik/ctx"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"google.golang.org/grpc/codes"
)

func (h *OrganizationHandler) OrganizationListGet(c context.Context, req *reqOrganizationCore.OrganizationListGet) (*OrganizationListGetRes, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("OrganizationListGet")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	param, err := req.ToParam(ctx, h.resolver)
	if err != nil {
		return nil, err
	}

	res, err := h.organizationService.OrganizationsGet(param)
	if err != nil {
		return nil, err
	}

	code := codes.OK

	return &OrganizationListGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(code),
			Message: code.String(),
		},
		Data: response.OrganizationListFromDTO(&res.Data),
		Meta: resPkg.ResponseMetaSetup(
			res.Total,
			param.Filter.Raw().LimitOrDefault(),
			param.Filter.Raw().PageOrDefault(),
		),
	}, nil
}
