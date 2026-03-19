// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"time"

	ctx "github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/groupperiod"
	trxParam "github.com/dedyf5/resik/entities/transaction/param"
	"github.com/dedyf5/resik/pkg/goku"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/datetime"
)

func (m *MerchantOmzetGet) ToParam(c *ctx.Ctx) (result *trxParam.MerchantOmzetGet, err *resPkg.Status) {
	orderStr := "period"
	if m.Order != nil {
		orderStr = m.GetOrder()
	}

	datetimeStart, err := datetime.FromString(m.GetDatetimeStart(), time.RFC3339, c)
	if err != nil {
		return nil, err
	}

	datetimeEnd, err := datetime.FromString(m.GetDatetimeEnd(), time.RFC3339, c)
	if err != nil {
		return nil, err
	}

	return &trxParam.MerchantOmzetGet{
		Ctx:        c,
		MerchantID: m.GetMerchantId(),
		GroupPeriod: groupperiod.GroupPeriod{
			Mode:          groupperiod.Mode(m.GetMode()),
			DatetimeStart: datetimeStart,
			DatetimeEnd:   datetimeEnd,
			Timezone:      m.GetTimezone(),
		},
		Filter: *goku.NewFilter(m.GetSearch(), m.GetPage(), m.GetLimit()),
		Orders: goku.OrdersBuilder(orderStr),
	}, nil
}
