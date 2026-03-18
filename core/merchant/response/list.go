// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"time"

	merchantEntity "github.com/dedyf5/resik/entities/merchant"
)

func MerchantListFromEntity(src []merchantEntity.Merchant) (res []*MerchantList) {
	for _, v := range src {
		res = append(res, &MerchantList{
			Id:        v.ID,
			Name:      v.Name,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Format(time.RFC3339),
		})
	}
	return
}
