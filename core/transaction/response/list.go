// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"time"

	dtoTrx "github.com/dedyf5/resik/core/transaction/dto"
)

func TransactionsGetFromDTO(src *dtoTrx.Transactions) (res []*TransactionList) {
	for _, v := range *src {
		res = append(res, &TransactionList{
			Id:           v.PublicID.String32(),
			BillTotal:    v.BillTotal,
			TransactedAt: v.TransactedAt.Format(time.RFC3339),
			Branch: &Branch{
				Id:   v.BranchPublicID.String32(),
				Name: v.BranchName,
				Organization: &Organization{
					Id:   v.OrganizationPublicID.String32(),
					Name: v.OrganizationName,
				},
			},
			Creator: UserFromUserEntity(v.Creator),
			Updater: UserFromUserEntity(v.Updater),
		})
	}
	return
}
