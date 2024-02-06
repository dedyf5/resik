// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import trxEntity "github.com/dedyf5/resik/entities/transaction"

type MerchantOmzet struct {
	MerchantName string  `json:"merchant_name" example:"Resik Merchant"`
	Omzet        float64 `json:"omzet" example:"1000000.5"`
	Period       string  `json:"period" example:"2024-02-06"`
}

func MerchantOmzetFromEntity(src []trxEntity.MerchantOmzet) []MerchantOmzet {
	res := make([]MerchantOmzet, 0, cap(src))
	for _, v := range src {
		res = append(res, MerchantOmzet{
			MerchantName: v.MerchantName,
			Omzet:        v.Omzet,
			Period:       v.Period,
		})
	}
	return res
}
