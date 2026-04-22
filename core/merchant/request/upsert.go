// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"time"

	"github.com/dedyf5/resik/ctx"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/datetime"
)

func (m *MerchantPost) ToEntity(ctx *ctx.Ctx) (res *merchantEntity.Merchant, err *resPkg.Status) {
	datetime, err := datetime.FromString(m.GetCreatedAt(), time.RFC3339, ctx)
	if err != nil {
		return nil, err
	}
	userID := ctx.UserClaims().UserID
	return &merchantEntity.Merchant{
		UserID:      userID,
		Name:        m.GetName(),
		Description: m.Description,
		CreatedBy:   userID,
		CreatedAt:   *datetime,
		UpdatedBy:   userID,
		UpdatedAt:   *datetime,
	}, nil
}

func (m *MerchantPut) ToEntity(ctx *ctx.Ctx) (res *merchantEntity.Merchant, err *resPkg.Status) {
	datetime, err := datetime.FromString(m.GetUpdatedAt(), time.RFC3339, ctx)
	if err != nil {
		return nil, err
	}
	return &merchantEntity.Merchant{
		ID:          m.GetId(),
		Name:        m.GetName(),
		Description: m.Description,
		UpdatedBy:   ctx.UserClaims().UserID,
		UpdatedAt:   *datetime,
	}, nil
}
