// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"context"

	"github.com/dedyf5/resik/app/grpc/proto/meta"
	"github.com/dedyf5/resik/app/grpc/proto/status"
	reqTrxCore "github.com/dedyf5/resik/core/transaction/req"
	resTrxCore "github.com/dedyf5/resik/core/transaction/res"
	"github.com/dedyf5/resik/ctx"
	"google.golang.org/grpc/codes"
)

func (r *OutletOmzetGetReq) ToCoreReq() *reqTrxCore.OutletOmzetGet {
	var page *uint = nil
	if r.Page != nil {
		tmp := uint(*r.Page)
		page = &tmp
	}
	var limit *uint = nil
	if r.Limit != nil {
		tmp := uint(*r.Limit)
		limit = &tmp
	}
	return &reqTrxCore.OutletOmzetGet{
		OutletID:      r.OutletID,
		Mode:          r.Mode,
		DateTimeStart: r.DateTimeStart,
		DateTimeEnd:   r.DateTimeEnd,
		Search:        r.Search,
		Order:         r.Order,
		Page:          page,
		Limit:         limit,
	}
}

func (h *TransactionHandler) OutletOmzetGet(c context.Context, req *OutletOmzetGetReq) (*OutletOmzetGetRes, error) {
	ctx, err := ctx.NewHTTPFromGRPC(c, h.log)
	if err != nil {
		return nil, err.GRPC().Err()
	}
	ctx.App.Logger().Debug("OutletOmzetGet")

	payload := req.ToCoreReq()
	if err := h.validator.Struct(payload, ctx.Lang); err != nil {
		return nil, err.GRPC().Err()
	}

	param := payload.ToParam(ctx)

	if _, err := ctx.UserClaims.OutletIDIsAccessible(param.OutletID); err != nil {
		return nil, err.GRPC().Err()
	}

	res, err := h.trxService.OutletOmzetGet(param)
	if err != nil {
		return nil, err.GRPC().Err()
	}

	return &OutletOmzetGetRes{
		Status: &status.Status{
			Code:    status.CodePlus(codes.OK),
			Message: codes.OK.String(),
		},
		Data: resTrxCore.OutletOmzetFromEntity(res.Data),
		Meta: meta.MetaSetup(res.Total, param.Filter.Limit, param.Filter.Page),
	}, nil
}
