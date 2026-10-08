// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"context"

	commonCore "github.com/dedyf5/resik/core/common"
	reqOrganizationCore "github.com/dedyf5/resik/core/organization/request"
	"github.com/dedyf5/resik/ctx"
	"google.golang.org/grpc/codes"
)

func (h *OrganizationHandler) OrganizationDelete(c context.Context, req *reqOrganizationCore.OrganizationDelete) (*commonCore.Empty, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("OrganizationDelete")

	if err = h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	organizationID, _, err := ctx.GetOrganizationID(h.resolver, req.GetId(), "organization:delete")
	if err != nil {
		return nil, err
	}

	if _, err = h.organizationService.OrganizationDelete(ctx, req.ToOrganization(organizationID)); err != nil {
		return nil, err
	}

	return &commonCore.Empty{
		Code:    commonCore.StatusCodePlus(codes.OK),
		Message: codes.OK.String(),
	}, nil
}
