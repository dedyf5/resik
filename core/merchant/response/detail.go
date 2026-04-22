// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"time"

	merchantEntity "github.com/dedyf5/resik/entities/merchant"
)

func MerchantDetailFromEntity(src *merchantEntity.Merchant) (res *MerchantDetail) {
	if src == nil {
		return nil
	}

	return &MerchantDetail{
		Id:     src.ID,
		UserId: src.UserID,
		User: &User{
			Id:       src.User.ID,
			Name:     src.User.Name,
			Username: src.User.Username,
		},
		Name:        src.Name,
		Description: src.Description,
		CreatedAt:   src.CreatedAt.Format(time.RFC3339),
		CreatedBy:   src.CreatedBy,
		UpdatedAt:   src.UpdatedAt.Format(time.RFC3339),
		UpdatedBy:   src.UpdatedBy,
	}
}
