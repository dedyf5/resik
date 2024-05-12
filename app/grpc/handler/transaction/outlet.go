// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/meta"
	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqTrxCore "github.com/dedyf5/resik/core/transaction/request"
	resTrxCore "github.com/dedyf5/resik/core/transaction/response"
	"github.com/dedyf5/resik/ctx"
	"google.golang.org/grpc/codes"
)

func (h *TransactionHandler) OutletOmzetGet(c context.Context, req *reqTrxCore.OutletOmzetGet) (*OutletOmzetGetRes, error) {
	ctx, err := ctx.NewHTTPFromGRPC(c, h.log)
	if err != nil {
		return nil, err.GRPC().Err()
	}
	ctx.App.Logger().Debug("OutletOmzetGet")

	if err := h.validator.Struct(req, ctx.Lang); err != nil {
		return nil, err.GRPC().Err()
	}

	param := req.ToParam(ctx)

	if _, err := ctx.UserClaims.OutletIDIsAccessible(param.OutletID); err != nil {
		return nil, err.GRPC().Err()
	}

	res, err := h.trxService.OutletOmzetGet(param)
	if err != nil {
		return nil, err.GRPC().Err()
	}

	return &OutletOmzetGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: resTrxCore.OutletOmzetFromEntity(res.Data),
		Meta: meta.MetaSetup(res.Total, param.Filter.Limit, param.Filter.Page),
	}, nil
}
