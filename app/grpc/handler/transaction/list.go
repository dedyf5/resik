// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"context"

	commonCore "github.com/dedyf5/resik/core/common"
	reqTransactionCore "github.com/dedyf5/resik/core/transaction/request"

	"github.com/dedyf5/resik/core/transaction/response"
	"github.com/dedyf5/resik/ctx"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"google.golang.org/grpc/codes"
)

func (h *TransactionHandler) TransactionListGet(c context.Context, req *reqTransactionCore.TransactionListGet) (*TransactionListGetRes, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("TransactionListGet")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	param, err := req.ToParam(ctx, h.resolver)
	if err != nil {
		return nil, err
	}

	res, err := h.trxService.TransactionsGet(param)
	if err != nil {
		return nil, err
	}

	code := codes.OK

	return &TransactionListGetRes{
		Status: &commonCore.Status{
			Code:    commonCore.StatusCodePlus(code),
			Message: code.String(),
		},
		Data: response.TransactionsGetFromDTO(&res.Data),
		Meta: resPkg.ResponseMetaSetup(
			res.Total,
			param.Filter.Raw().LimitOrDefault(),
			param.Filter.Raw().PageOrDefault(),
		),
	}, nil
}
