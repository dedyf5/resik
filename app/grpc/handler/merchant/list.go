// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqMerchantCore "github.com/dedyf5/resik/core/merchant/request"
	"github.com/dedyf5/resik/core/merchant/response"
	"github.com/dedyf5/resik/ctx"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"google.golang.org/grpc/codes"
)

func (h *MerchantHandler) MerchantListGet(c context.Context, req *reqMerchantCore.MerchantListGet) (*MerchantListGetRes, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("MerchantListGet")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	param := req.ToParam(ctx)

	res, err := h.merchantService.MerchantListGet(param)
	if err != nil {
		return nil, err
	}

	code := codes.OK

	return &MerchantListGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(code),
			Message: code.String(),
		},
		Data: response.MerchantListFromEntity(res.Data),
		Meta: resPkg.ResponseMetaSetup(
			res.Total,
			param.Filter.Raw().LimitOrDefault(),
			param.Filter.Raw().PageOrDefault(),
		),
	}, nil
}
