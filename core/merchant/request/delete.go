// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import "github.com/dedyf5/resik/entities/merchant"

func (*MerchantDelete) ToMerchant(merchantID uint64) *merchant.Merchant {
	return &merchant.Merchant{
		ID: merchantID,
	}
}
