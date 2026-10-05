// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"time"

	commonCore "github.com/dedyf5/resik/core/common"
	dtoTrx "github.com/dedyf5/resik/core/transaction/dto"
)

func TransactionDetailFromDTO(src *dtoTrx.Transaction) *TransactionDetail {
	if src == nil {
		return nil
	}

	return &TransactionDetail{
		Id:           src.PublicID.String32(),
		BillTotal:    src.BillTotal,
		TransactedAt: src.TransactedAt.Format(time.RFC3339),
		Branch: &Branch{
			Id:   src.BranchPublicID.String32(),
			Name: src.BranchName,
			Organization: &Organization{
				Id:   src.OrganizationPublicID.String32(),
				Name: src.OrganizationName,
			},
		},
		CreatedAt: src.CreatedAt.Format(time.RFC3339),
		Creator: &commonCore.User{
			Id:       src.Creator.PublicID.String32(),
			Name:     src.Creator.Name,
			Username: src.Creator.Username,
		},
		UpdatedAt: src.UpdatedAt.Format(time.RFC3339),
		Updater: &commonCore.User{
			Id:       src.Updater.PublicID.String32(),
			Name:     src.Updater.Name,
			Username: src.Updater.Username,
		},
	}
}
