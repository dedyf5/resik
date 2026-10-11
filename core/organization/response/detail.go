// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"time"

	commonCore "github.com/dedyf5/resik/core/common"
	dtoOrganization "github.com/dedyf5/resik/core/organization/dto"
)

func OrganizationDetailFromDTO(src *dtoOrganization.Organization) *OrganizationDetail {
	if src == nil {
		return nil
	}

	return &OrganizationDetail{
		Id:          src.PublicID.String32(),
		Name:        src.Name,
		Description: src.Description,
		CreatedAt:   src.CreatedAt.Format(time.RFC3339),
		Creator:     commonCore.UserFromUserEntity(src.Creator),
		UpdatedAt:   src.UpdatedAt.Format(time.RFC3339),
		Updater:     commonCore.UserFromUserEntity(src.Updater),
	}
}
