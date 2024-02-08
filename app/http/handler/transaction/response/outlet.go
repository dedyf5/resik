// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import trxEntity "github.com/dedyf5/resik/entities/transaction"

type OutletOmzet struct {
	MerchantName string  `json:"merchant_name" example:"Resik Merchant"`
	OutletName   string  `json:"outlet_name" example:"Resik Outlet"`
	Omzet        float64 `json:"omzet" example:"50000.5"`
	Period       string  `json:"period" example:"2024-02-06"`
}

func OutletOmzetFromEntity(src []trxEntity.OutletOmzet) []OutletOmzet {
	res := make([]OutletOmzet, 0, cap(src))
	for _, v := range src {
		res = append(res, OutletOmzet{
			MerchantName: v.MerchantName,
			OutletName:   v.OutletName,
			Omzet:        v.Omzet,
			Period:       v.Period,
		})
	}
	return res
}
