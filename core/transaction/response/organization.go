// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import trxEntity "github.com/dedyf5/resik/entities/transaction"

func OrganizationOmzetFromEntity(src []trxEntity.OrganizationOmzet) []*OrganizationOmzet {
	res := make([]*OrganizationOmzet, 0, cap(src))
	for _, v := range src {
		res = append(res, &OrganizationOmzet{
			OrganizationName: v.OrganizationName,
			Omzet:            v.Omzet,
			Period:           v.Period,
		})
	}
	return res
}
