// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"github.com/dedyf5/resik/entities/merchant"
)

type MerchantDelete struct {
	ID uint64 `json:"-" param:"id" query:"-" validate:"required" example:"123"`
}

func (m *MerchantDelete) ToMerchant() *merchant.Merchant {
	return &merchant.Merchant{
		ID: m.ID,
	}
}
