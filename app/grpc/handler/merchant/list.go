// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/meta"
	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqMerchantCore "github.com/dedyf5/resik/core/merchant/request"
	"github.com/dedyf5/resik/core/merchant/response"
	"github.com/dedyf5/resik/ctx"
	"google.golang.org/grpc/codes"
)

func (h *MerchantHandler) MerchantListGet(c context.Context, req *reqMerchantCore.MerchantListGet) (*MerchantListGetRes, error) {
	ctx, err := ctx.NewHTTPFromGRPC(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.App.Logger().Debug("MerchantListGet")

	if err := h.validator.Struct(req, ctx.Lang); err != nil {
		return nil, err
	}

	param := req.ToParam(ctx)

	res, err := h.merchantService.MerchantListGet(param)
	if err != nil {
		return nil, err
	}

	code := codes.OK
	if len(res.Data) == 0 {
		code = codes.NotFound
	}

	return &MerchantListGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(code),
			Message: code.String(),
		},
		Data: response.MerchantListFromEntity(res.Data),
		Meta: meta.MetaSetup(res.Total, param.Filter.Limit, param.Filter.Page),
	}, nil
}
