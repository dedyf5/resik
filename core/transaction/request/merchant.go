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

func (m *MerchantOmzetGet) ToParam(c *ctx.Ctx) *trxParam.MerchantOmzetGet {
	orderStr := "period"
	if m.Order != nil {
		orderStr = *m.Order
	}
	var page int = 1
	if m.Page != nil {
		page = int(*m.Page)
	}
	var limit int = 10
	if m.Limit != nil {
		limit = int(*m.Limit)
	}
	return &trxParam.MerchantOmzetGet{
		Ctx:        c,
		MerchantID: m.MerchantID,
		GroupPeriod: groupperiod.GroupPeriod{
			Mode:          groupperiod.Mode(m.Mode),
			DatetimeStart: m.DateTimeStart,
			DatetimeEnd:   m.DateTimeEnd,
		},
		Filter: goku.Filter{
			Search: m.Search,
			Page:   page,
			Limit:  limit,
		},
		Orders: goku.OrdersBuilder(orderStr),
	}
}
