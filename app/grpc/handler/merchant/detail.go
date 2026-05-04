// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"context"
	"net/http"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqMerchantCore "github.com/dedyf5/resik/core/merchant/request"
	resMerchantCore "github.com/dedyf5/resik/core/merchant/response"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"google.golang.org/grpc/codes"
)

func (h *MerchantHandler) MerchantDetailGet(c context.Context, req *reqMerchantCore.MerchantDetailGet) (*MerchantDetailGetRes, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("MerchantDetailGet")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	merchantID, err := ctx.UserClaims().GetMerchantID(req.GetId())
	if err != nil {
		return nil, err
	}

	merchant, err := h.merchantService.MerchantGetByIDAndOwnerID(ctx, merchantID, ctx.UserClaims().UserID())
	if err != nil {
		return nil, err
	}

	if merchant == nil {
		localizer := ctx.Lang().Localizer
		return nil, resPkg.NewStatusMessage(
			http.StatusNotFound,
			term.NotFoundVal.Localize(localizer, term.Merchant.Localize(localizer)),
			nil,
		)
	}

	return &MerchantDetailGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: resMerchantCore.MerchantDetailFromDTO(merchant),
	}, nil
}
