// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	commonCore "github.com/dedyf5/resik/core/common"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
)

func BranchOmzetFromEntity(src []trxEntity.BranchOmzet) []*BranchOmzet {
	res := make([]*BranchOmzet, 0, cap(src))
	for _, v := range src {
		res = append(res, &BranchOmzet{
			Branch: &commonCore.Branch{
				Id:   v.BranchPublicID.String32(),
				Name: v.BranchName,
				Organization: &commonCore.Organization{
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
