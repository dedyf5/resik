// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqMerchantCore "github.com/dedyf5/resik/core/merchant/request"
	"github.com/dedyf5/resik/ctx"
	"google.golang.org/grpc/codes"
)

func (h *MerchantHandler) MerchantDelete(c context.Context, req *reqMerchantCore.MerchantDelete) (*status.Empty, error) {
	ctx, err := ctx.NewHTTPFromGRPC(c, h.log)
	if err != nil {
		return nil, err.GRPC().Err()
	}
	ctx.App.Logger().Debug("MerchantDelete")

	if err = h.validator.Struct(req, ctx.Lang); err != nil {
		return nil, err.GRPC().Err()
	}

	if _, err = ctx.UserClaims.MerchantIDIsAccessible(req.ID); err != nil {
		return nil, err.GRPC().Err()
	}

	if _, err = h.merchantService.MerchantDelete(ctx, req.ToMerchant()); err != nil {
		return nil, err.GRPC().Err()
	}

	return &status.Empty{
		Code:    status.CodePlus(codes.OK),
		Message: codes.OK.String(),
	}, nil
}
