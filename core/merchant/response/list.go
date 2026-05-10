// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	dtoMerchant "github.com/dedyf5/resik/core/merchant/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MerchantListFromDTO(src *dtoMerchant.Merchants) (res []*MerchantList) {
	for _, v := range *src {
		res = append(res, &MerchantList{
			Id:        v.PublicID.String32(),
			Name:      v.Name,
			CreatedAt: timestamppb.New(v.CreatedAt),
			UpdatedAt: timestamppb.New(v.UpdatedAt),
		})
	}
	return
}
