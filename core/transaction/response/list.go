// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"time"

	commonCore "github.com/dedyf5/resik/core/common"
	dtoTrx "github.com/dedyf5/resik/core/transaction/dto"
)

func TransactionsGetFromDTO(src *dtoTrx.Transactions) (res []*TransactionList) {
	for _, v := range *src {
		res = append(res, &TransactionList{
			Id:           v.PublicID.String32(),
			BillTotal:    v.BillTotal,
			TransactedAt: v.TransactedAt.Format(time.RFC3339),
			Branch: &commonCore.Branch{
				Id:   v.BranchPublicID.String32(),
				Name: v.BranchName,
				Organization: &commonCore.Organization{
					Id:   v.OrganizationPublicID.String32(),
					Name: v.OrganizationName,
				},
			},
			Creator: commonCore.UserFromUserEntity(v.Creator),
			Updater: commonCore.UserFromUserEntity(v.Updater),
		})
	}
	return
}
