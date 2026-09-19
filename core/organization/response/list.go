// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"time"

	dtoOrganization "github.com/dedyf5/resik/core/organization/dto"
)

func OrganizationListFromDTO(src *dtoOrganization.Organizations) (res []*OrganizationList) {
	for _, v := range *src {
		res = append(res, &OrganizationList{
			Id:        v.PublicID.String32(),
			Name:      v.Name,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Format(time.RFC3339),
		})
	}
	return
}
