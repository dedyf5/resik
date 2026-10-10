// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"context"

	commonCore "github.com/dedyf5/resik/core/common"
	dtoTrx "github.com/dedyf5/resik/core/transaction/dto"
	reqTrxCore "github.com/dedyf5/resik/core/transaction/request"
	"github.com/dedyf5/resik/ctx"
	"google.golang.org/grpc/codes"
)

func (h *TransactionHandler) TransactionDelete(c context.Context, req *reqTrxCore.Id) (*commonCore.Empty, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("TransactionDelete")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	param, err := req.ToParam(ctx, h.resolver, dtoTrx.PERMISSION_CODE_DELETE)
	if err != nil {
		return nil, err
	}

	if ok, err := h.trxService.TransactionDelete(param); err != nil || !ok {
		return nil, err
	}

	return &commonCore.Empty{
		Code:    commonCore.StatusCodePlus(codes.OK),
		Message: codes.OK.String(),
	}, nil
}
