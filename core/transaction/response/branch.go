// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import trxEntity "github.com/dedyf5/resik/entities/transaction"

func BranchOmzetFromEntity(src []trxEntity.BranchOmzet) []*BranchOmzet {
	res := make([]*BranchOmzet, 0, cap(src))
	for _, v := range src {
		res = append(res, &BranchOmzet{
			OrganizationName: v.OrganizationName,
			BranchName:       v.BranchName,
			Omzet:            v.Omzet,
			Period:           v.Period,
		})
	}
	return res
}
