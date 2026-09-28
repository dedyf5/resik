// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import trxEntity "github.com/dedyf5/resik/entities/transaction"

func BranchOmzetFromEntity(src []trxEntity.BranchOmzet) []*BranchOmzet {
	res := make([]*BranchOmzet, 0, cap(src))
	for _, v := range src {
		res = append(res, &BranchOmzet{
			Branch: &Branch{
				Id:   v.BranchPublicID.String32(),
				Name: v.BranchName,
				Organization: &Organization{
					Id:   v.OrganizationPublicID.String32(),
					Name: v.OrganizationName,
				},
			},
			Omzet:  v.Omzet,
			Period: v.Period,
		})
	}
	return res
}
