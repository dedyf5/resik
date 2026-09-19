// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqOrganizationCore "github.com/dedyf5/resik/core/organization/request"
	"github.com/dedyf5/resik/ctx"
	"google.golang.org/grpc/codes"
)

func (h *OrganizationHandler) OrganizationDelete(c context.Context, req *reqOrganizationCore.OrganizationDelete) (*status.Empty, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("OrganizationDelete")

	if err = h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	organizationID, err := ctx.GetOrganizationID(h.resolver, req.GetId(), "organization:delete")
	if err != nil {
		return nil, err
	}

	if _, err = h.organizationService.OrganizationDelete(ctx, req.ToOrganization(organizationID)); err != nil {
		return nil, err
	}

	return &status.Empty{
		Code:    status.CodePlus(codes.OK),
		Message: codes.OK.String(),
	}, nil
}
