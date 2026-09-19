// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqTrxCore "github.com/dedyf5/resik/core/transaction/request"
	resTrxCore "github.com/dedyf5/resik/core/transaction/response"
	"github.com/dedyf5/resik/ctx"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"google.golang.org/grpc/codes"
)

func (h *TransactionHandler) OrganizationOmzetGet(c context.Context, req *reqTrxCore.OrganizationOmzetGet) (*OrganizationOmzetGetRes, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("OrganizationOmzetGet")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	organizationID, err := ctx.GetOrganizationID(h.resolver, req.GetOrganizationId(), "transaction:read")
	if err != nil {
		return nil, err
	}

	param, err := req.ToParam(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	res, err := h.trxService.OrganizationOmzetGet(param)
	if err != nil {
		return nil, err
	}

	return &OrganizationOmzetGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: resTrxCore.OrganizationOmzetFromEntity(res.Data),
		Meta: resPkg.ResponseMetaSetup(
			res.Total,
			param.Filter.Raw().LimitOrDefault(),
			param.Filter.Raw().PageOrDefault(),
		),
	}, nil
}
