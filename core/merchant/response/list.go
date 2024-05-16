// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	"github.com/dedyf5/resik/utils/datetime"
)

func MerchantListFromEntity(src []merchantEntity.Merchant) (res []*MerchantList) {
	format := datetime.FormatyyyyMMddHHmmss.ToString()
	for _, v := range src {
		res = append(res, &MerchantList{
			ID:        v.ID,
			Name:      v.Name,
			CreatedAt: v.CreatedAt.Format(format),
			UpdatedAt: v.UpdatedAt.Format(format),
		})
	}
	return
}
