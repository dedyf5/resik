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

func (h *TransactionHandler) MerchantOmzetGet(c context.Context, req *reqTrxCore.MerchantOmzetGet) (*MerchantOmzetGetRes, error) {
	ctx, err := ctx.NewHTTPFromGRPC(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.App.Logger().Debug("MerchantOmzetGet")

	if err := h.validator.Struct(req, ctx.Lang); err != nil {
		return nil, err
	}

	param := req.ToParam(ctx)

	if _, err := ctx.UserClaims.MerchantIDIsAccessible(param.MerchantID); err != nil {
		return nil, err
	}

	res, err := h.trxService.MerchantOmzetGet(param)
	if err != nil {
		return nil, err
	}

	return &MerchantOmzetGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: resTrxCore.MerchantOmzetFromEntity(res.Data),
		Meta: meta.MetaSetup(res.Total, param.Filter.Limit, param.Filter.Page),
	}, nil
}
