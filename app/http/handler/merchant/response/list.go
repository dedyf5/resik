// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	"github.com/dedyf5/resik/utils/datetime"
)

type MerchantList struct {
	ID        uint64 `json:"id" example:"123"`
	Name      string `json:"name" example:"Resik"`
	CreatedAt string `json:"created_at" example:"2024-01-14 11:40:00"`
	UpdatedAt string `json:"updated_at" example:"2024-01-14 11:40:00"`
}

func MerchantListFromEntity(src []merchantEntity.Merchant) (res []MerchantList) {
	format := datetime.FormatyyyyMMddHHmmss.ToString()
	for _, v := range src {
		res = append(res, MerchantList{
			ID:        v.ID,
			Name:      v.Name,
			CreatedAt: v.CreatedAt.Format(format),
			UpdatedAt: v.UpdatedAt.Format(format),
		})
	}
	return
}
