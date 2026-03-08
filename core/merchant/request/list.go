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
	page := 1
	if m.Page != nil {
		page = int(m.GetPage())
	}
	limit := 10
	if m.Limit != nil {
		limit = int(m.GetLimit())
	}
	return &paramMerchant.MerchantListGet{
		Ctx:         c,
		MerchantIDs: c.UserClaims().MerchantIDs,
		Filter: goku.Filter{
			Search: m.GetSearch(),
			Page:   page,
			Limit:  limit,
		},
		Orders: goku.OrdersBuilder(orderStr),
	}
}
