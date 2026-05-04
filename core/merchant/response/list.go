// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"time"

	dtoMerchant "github.com/dedyf5/resik/core/merchant/dto"
)

func MerchantListFromDTO(src *dtoMerchant.Merchants) (res []*MerchantList) {
	for _, v := range *src {
		res = append(res, &MerchantList{
			Id:        v.PublicID.String32(),
			Name:      v.Name,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Format(time.RFC3339),
		})
	}
	return
}
