// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"context"

	commonCore "github.com/dedyf5/resik/core/common"
	reqTrxCore "github.com/dedyf5/resik/core/transaction/request"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	"google.golang.org/grpc/codes"
)

func (h *TransactionHandler) TransactionPost(c context.Context, req *reqTrxCore.TransactionPost) (*commonCore.StatusAndId, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("TransactionPost")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	entity, err := req.ToEntity(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := h.trxService.TransactionInsert(ctx, entity); err != nil {
		return nil, err
	}

	return &commonCore.StatusAndId{
		Status: &commonCore.Status{
			Code: commonCore.StatusCodePlus(codes.OK),
			Message: term.SuccessfullyCreatedVal.Localize(
				ctx.Lang().Localizer,
				term.Transaction.Localize(ctx.Lang().Localizer),
			),
		},
		Data: &commonCore.Id{
			Id: entity.PublicID.String32(),
		},
	}, nil
}

func (h *TransactionHandler) TransactionPut(c context.Context, req *reqTrxCore.TransactionPut) (*commonCore.StatusAndId, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("TransactionPut")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	entity, err := req.ToEntity(ctx)
	if err != nil {
		return nil, err
	}

	_, err = h.trxService.TransactionUpdate(ctx, entity)
	if err != nil {
		return nil, err
	}

	return &commonCore.StatusAndId{
		Status: &commonCore.Status{
			Code: commonCore.StatusCodePlus(codes.OK),
			Message: term.SuccessfullyUpdatedVal.Localize(
				ctx.Lang().Localizer,
				term.Transaction.Localize(ctx.Lang().Localizer),
			),
		},
		Data: &commonCore.Id{
			Id: req.GetId(),
		},
	}, nil
}
