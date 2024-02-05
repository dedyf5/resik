package dto

import trxEntity "github.com/dedyf5/resik/entities/transaction"

type MerchantOmzet struct {
	Data  []trxEntity.MerchantOmzet
	Total uint64
}
