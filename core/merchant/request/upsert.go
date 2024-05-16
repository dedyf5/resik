// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"github.com/dedyf5/resik/ctx"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/datetime"
)

func (m *MerchantPost) ToEntity(ctx *ctx.Ctx) (res *merchantEntity.Merchant, status *resPkg.Status) {
	datetime, err := datetime.FromString(m.CreatedAt, datetime.FormatyyyyMMddHHmmss)
	if err != nil {
		return nil, err
	}
	return &merchantEntity.Merchant{
		UserID:    ctx.UserClaims.UserID,
		Name:      m.Name,
		CreatedBy: ctx.UserClaims.UserID,
		CreatedAt: *datetime,
		UpdatedBy: ctx.UserClaims.UserID,
		UpdatedAt: *datetime,
	}, nil
}
