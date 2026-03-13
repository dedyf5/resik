// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"github.com/dedyf5/resik/ctx"
	paramMerchant "github.com/dedyf5/resik/entities/merchant/param"
	"github.com/dedyf5/resik/pkg/goku"
)

func (m *MerchantListGet) ToParam(c *ctx.Ctx) *paramMerchant.MerchantListGet {
	orderStr := "name"
	if m.Order != nil {
		orderStr = m.GetOrder()
	}
	return &paramMerchant.MerchantListGet{
		Ctx:         c,
		MerchantIDs: c.UserClaims().MerchantIDs,
		Filter:      *goku.NewFilter(m.GetSearch(), m.GetPage(), m.GetLimit()),
		Orders:      goku.OrdersBuilder(orderStr),
	}
}
