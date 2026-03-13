// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package dto

import trxEntity "github.com/dedyf5/resik/entities/transaction"

type MerchantOmzet struct {
	Data  []trxEntity.MerchantOmzet
	Total int64
}
