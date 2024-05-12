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
		orderStr = *o.Order
	}
	var page int = 1
	if o.Page != nil {
		page = int(*o.Page)
	}
	var limit int = 10
	if o.Limit != nil {
		limit = int(*o.Limit)
	}
	return &trxParam.OutletOmzetGet{
		Ctx:      c,
		OutletID: o.OutletID,
		GroupPeriod: groupperiod.GroupPeriod{
			Mode:          groupperiod.Mode(o.Mode),
			DatetimeStart: o.DateTimeStart,
			DatetimeEnd:   o.DateTimeEnd,
		},
		Filter: goku.Filter{
			Search: o.Search,
			Page:   page,
			Limit:  limit,
		},
		Orders: goku.OrdersBuilder(orderStr),
	}
}
