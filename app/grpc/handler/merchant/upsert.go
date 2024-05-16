// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqMerchantCore "github.com/dedyf5/resik/core/merchant/request"
	resMerchantCore "github.com/dedyf5/resik/core/merchant/response"
	"github.com/dedyf5/resik/ctx"
	"google.golang.org/grpc/codes"
)

func (h *MerchantHandler) MerchantPost(c context.Context, req *reqMerchantCore.MerchantPost) (*MerchantUpsertRes, error) {
	ctx, err := ctx.NewHTTPFromGRPC(c, h.log)
	if err != nil {
		return nil, err.GRPC().Err()
	}
	ctx.App.Logger().Debug("MerchantPost")

	if err := h.validator.Struct(req, ctx.Lang); err != nil {
		return nil, err.GRPC().Err()
	}

	entity, err := req.ToEntity(ctx)
	if err != nil {
		return nil, err.GRPC().Err()
	}

	if _, err := h.merchantService.MerchantInsert(ctx, entity); err != nil {
		return nil, err.GRPC().Err()
	}

	return &MerchantUpsertRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: &resMerchantCore.MerchantUpsert{
			ID: entity.ID,
		},
	}, nil
}

func (h *MerchantHandler) MerchantPut(c context.Context, req *reqMerchantCore.MerchantPut) (*MerchantUpsertRes, error) {
	ctx, err := ctx.NewHTTPFromGRPC(c, h.log)
	if err != nil {
		return nil, err.GRPC().Err()
	}
	ctx.App.Logger().Debug("MerchantPut")

	if err := h.validator.Struct(req, ctx.Lang); err != nil {
		return nil, err.GRPC().Err()
	}

	entity, err := req.ToEntity(ctx)
	if err != nil {
		return nil, err.GRPC().Err()
	}

	if _, err = ctx.UserClaims.MerchantIDIsAccessible(entity.ID); err != nil {
		return nil, err.GRPC().Err()
	}

	if _, err = h.merchantService.MerchantUpdate(ctx, entity); err != nil {
		return nil, err.GRPC().Err()
	}

	return &MerchantUpsertRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: &resMerchantCore.MerchantUpsert{
			ID: entity.ID,
		},
	}, nil
}
