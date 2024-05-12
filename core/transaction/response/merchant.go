// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import trxEntity "github.com/dedyf5/resik/entities/transaction"

func MerchantOmzetFromEntity(src []trxEntity.MerchantOmzet) []*MerchantOmzet {
	res := make([]*MerchantOmzet, 0, cap(src))
	for _, v := range src {
		res = append(res, &MerchantOmzet{
			MerchantName: v.MerchantName,
			Omzet:        v.Omzet,
			Period:       v.Period,
		})
	}
	return res
}
