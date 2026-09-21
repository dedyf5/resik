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

func (h *TransactionHandler) BranchOmzetGet(c context.Context, req *reqTrxCore.BranchOmzetGet) (*BranchOmzetGetRes, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("BranchOmzetGet")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	_, branchPublicID, err := ctx.GetBranchID(h.resolver, req.GetBranchId(), "transaction:read")
	if err != nil {
		return nil, err
	}

	param, err := req.ToParam(ctx, branchPublicID)
	if err != nil {
		return nil, err
	}

	res, err := h.trxService.BranchOmzetGet(param)
	if err != nil {
		return nil, err
	}

	return &BranchOmzetGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: resTrxCore.BranchOmzetFromEntity(res.Data),
		Meta: resPkg.ResponseMetaSetup(
			res.Total,
			param.Filter.Raw().LimitOrDefault(),
			param.Filter.Raw().PageOrDefault(),
		),
	}, nil
}
