// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"time"

	"github.com/dedyf5/resik/ctx"
	organizationEntity "github.com/dedyf5/resik/entities/organization"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/datetime"
)

func (m *OrganizationPost) ToEntity(ctx *ctx.Ctx) (res *organizationEntity.Organization, err *resPkg.Status) {
	datetime, err := datetime.FromString(m.GetCreatedAt(), time.RFC3339, ctx)
	if err != nil {
		return nil, err
	}
	userPublicID := ctx.UserClaims().UserPublicID()
	return &organizationEntity.Organization{
		Name:              m.GetName(),
		Description:       m.Description,
		CreatedByPublicID: userPublicID,
		CreatedAt:         *datetime,
		UpdatedByPublicID: userPublicID,
		UpdatedAt:         *datetime,
	}, nil
}

func (m *OrganizationPut) ToEntity(ctx *ctx.Ctx, organizationID uint64) (res *organizationEntity.Organization, err *resPkg.Status) {
	datetime, err := datetime.FromString(m.GetUpdatedAt(), time.RFC3339, ctx)
	if err != nil {
		return nil, err
	}

	return &organizationEntity.Organization{
		ID:                organizationID,
		Name:              m.GetName(),
		Description:       m.Description,
		UpdatedByPublicID: ctx.UserClaims().UserPublicID(),
		UpdatedAt:         *datetime,
	}, nil
}
