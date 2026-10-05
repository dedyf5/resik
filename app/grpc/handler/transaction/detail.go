// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"context"
	"net/http"

	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqTrxCore "github.com/dedyf5/resik/core/transaction/request"
	resTrxCore "github.com/dedyf5/resik/core/transaction/response"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"google.golang.org/grpc/codes"
)

func (h *TransactionHandler) TransactionDetailGet(c context.Context, req *reqTrxCore.TransactionDetailGet) (*TransactionDetailRes, error) {
	ctx, err := ctx.NewCtx(c, h.log)
	if err != nil {
		return nil, err
	}
	ctx.Log().Debug("TransactionDetailGet")

	if err := h.validator.Struct(req, ctx.Lang()); err != nil {
		return nil, err
	}

	param, err := req.ToParam(ctx, h.resolver)
	if err != nil {
		return nil, err
	}

	tx, err := h.trxService.TransactionGetByPublicID(param)
	if err != nil {
		return nil, err
	}

	if tx == nil {
		localizer := ctx.Lang().Localizer
		return nil, resPkg.NewStatusMessage(
			http.StatusNotFound,
			term.NotFoundVal.Localize(localizer, term.Transaction.Localize(localizer)),
			nil,
		)
	}

	return &TransactionDetailRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: resTrxCore.TransactionDetailFromDTO(tx),
	}, nil
}
