// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/jwt"
	paramMerchant "github.com/dedyf5/resik/entities/merchant/param"
	"github.com/dedyf5/resik/pkg/collection"
	"github.com/dedyf5/resik/pkg/goku"
)

func (m *MerchantListGet) ToParam(c *ctx.Ctx) *paramMerchant.MerchantsGet {
	orderStr := "name"
	if m.Order != nil {
		orderStr = m.GetOrder()
	}

	merchantIDs := collection.Map(c.UserClaims().Merchants, func(v jwt.Base) uint64 {
		return v.ID
	})

	return &paramMerchant.MerchantsGet{
		Ctx:         c,
		MerchantIDs: merchantIDs,
		Filter:      *goku.NewFilter(m.GetSearch(), m.GetPage(), m.GetLimit()),
		Orders:      goku.OrdersBuilder(orderStr),
	}
}
