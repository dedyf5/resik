// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"github.com/dedyf5/resik/ctx"
	commonEntity "github.com/dedyf5/resik/entities/common"
	paramMerchant "github.com/dedyf5/resik/entities/merchant/param"
	"github.com/dedyf5/resik/pkg/goku"
)

type MerchantListGet struct {
	commonEntity.Request
	Search string  `json:"search" query:"search" validate:"" example:"keyword"`
	Order  *string `json:"order" query:"order" validate:"omitempty,oneof_order=name created_at updated_at" example:"name"`
	Page   *uint   `json:"page" query:"page" validate:"omitempty,min=1" example:"1"`
	Limit  *uint   `json:"limit" query:"limit" validate:"omitempty,min=1,max=100" example:"10"`
}

func (m *MerchantListGet) ToParam(c *ctx.Ctx) *paramMerchant.MerchantListGet {
	orderStr := "name"
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
	return &paramMerchant.MerchantListGet{
		Ctx:         c,
		MerchantIDs: c.UserClaims.MerchantIDs,
		Filter: goku.Filter{
			Search: m.Search,
			Page:   page,
			Limit:  limit,
		},
		Orders: goku.OrdersBuilder(orderStr),
	}
}
