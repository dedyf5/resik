// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import trxEntity "github.com/dedyf5/resik/entities/transaction"

type MerchantOmzet struct {
	MerchantName string  `json:"merchant_name"`
	Omzet        float64 `json:"omzet"`
	Period       string  `json:"period"`
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
