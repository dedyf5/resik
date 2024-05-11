// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package res

import trxEntity "github.com/dedyf5/resik/entities/transaction"

func OutletOmzetFromEntity(src []trxEntity.OutletOmzet) []*OutletOmzet {
	res := make([]*OutletOmzet, 0, cap(src))
	for _, v := range src {
		res = append(res, &OutletOmzet{
			MerchantName: v.MerchantName,
			OutletName:   v.OutletName,
			Omzet:        v.Omzet,
			Period:       v.Period,
		})
	}
	return res
}
