// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package dto

import "github.com/dedyf5/resik/entities/merchant"

type MerchantList struct {
	Data  []merchant.Merchant
	Total uint64
}
