// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"time"

	dtoMerchant "github.com/dedyf5/resik/core/merchant/dto"
)

func MerchantDetailFromDTO(src *dtoMerchant.Merchant) *MerchantDetail {
	if src == nil {
		return nil
	}

	return &MerchantDetail{
		Id: src.PublicID.String32(),
		Owner: &User{
			Id:       src.Owner.ID,
			Name:     src.Owner.Name,
			Username: src.Owner.Username,
		},
		Name:        src.Name,
		Description: src.Description,
		CreatedAt:   src.CreatedAt.Format(time.RFC3339),
		Creator: &User{
			Id:       src.Creator.ID,
			Name:     src.Creator.Name,
			Username: src.Creator.Username,
		},
		UpdatedAt: src.UpdatedAt.Format(time.RFC3339),
		Updater: &User{
			Id:       src.Updater.ID,
			Name:     src.Updater.Name,
			Username: src.Updater.Username,
		},
	}
}
