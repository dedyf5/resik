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
		orderStr = m.GetOrder()
	}
	return &trxParam.MerchantOmzetGet{
		Ctx:        c,
		MerchantID: m.GetMerchantId(),
		GroupPeriod: groupperiod.GroupPeriod{
			Mode:          groupperiod.Mode(m.GetMode()),
			DatetimeStart: m.GetDatetimeStart(),
			DatetimeEnd:   m.GetDatetimeEnd(),
		},
		Filter: *goku.NewFilter(m.GetSearch(), m.GetPage(), m.GetLimit()),
		Orders: goku.OrdersBuilder(orderStr),
	}
}
