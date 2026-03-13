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
	return &trxParam.OutletOmzetGet{
		Ctx:      c,
		OutletID: o.GetOutletId(),
		GroupPeriod: groupperiod.GroupPeriod{
			Mode:          groupperiod.Mode(o.GetMode()),
			DatetimeStart: o.GetDatetimeStart(),
			DatetimeEnd:   o.GetDatetimeEnd(),
		},
		Filter: *goku.NewFilter(o.GetSearch(), o.GetPage(), o.GetLimit()),
		Orders: goku.OrdersBuilder(orderStr),
	}
}
