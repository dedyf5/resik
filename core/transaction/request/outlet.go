// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	ctx "github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/groupperiod"
	trxParam "github.com/dedyf5/resik/entities/transaction/param"
	"github.com/dedyf5/resik/pkg/goku"
)

func (o *OutletOmzetGet) ToParam(c *ctx.Ctx) *trxParam.OutletOmzetGet {
	orderStr := "period"
	if o.Order != nil {
		orderStr = o.GetOrder()
	}
	page := 1
	if o.Page != nil {
		page = int(o.GetPage())
	}
	limit := 10
	if o.Limit != nil {
		limit = int(o.GetLimit())
	}
	return &trxParam.OutletOmzetGet{
		Ctx:      c,
		OutletID: o.GetOutletId(),
		GroupPeriod: groupperiod.GroupPeriod{
			Mode:          groupperiod.Mode(o.GetMode()),
			DatetimeStart: o.GetDatetimeStart(),
			DatetimeEnd:   o.GetDatetimeEnd(),
		},
		Filter: goku.Filter{
			Search: o.GetSearch(),
			Page:   page,
			Limit:  limit,
		},
		Orders: goku.OrdersBuilder(orderStr),
	}
}
