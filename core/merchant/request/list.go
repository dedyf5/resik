// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"net/http"

	"github.com/dedyf5/resik/ctx"
	paramMerchant "github.com/dedyf5/resik/entities/merchant/param"
	"github.com/dedyf5/resik/internal/identity"
	"github.com/dedyf5/resik/pkg/goku"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func (m *MerchantListGet) ToParam(c *ctx.Ctx, resolver identity.IdentityResolver) (result *paramMerchant.MerchantsGet, err *resPkg.Status) {
	orderStr := "name"
	if m.Order != nil {
		orderStr = m.GetOrder()
	}

	merchantIDs, errRes := resolver.GetMerchantIDsByPermission(c.Context, c.UserClaims().UserID(), "merchant:read")
	if errRes != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errRes)
	}

	return &paramMerchant.MerchantsGet{
		Ctx:         c,
		MerchantIDs: merchantIDs,
		Filter:      *goku.NewFilter(m.GetSearch(), m.GetPage(), m.GetLimit()),
		Orders:      goku.OrdersBuilder(orderStr),
	}, nil
}
