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

type MerchantPost struct {
	Name      string `json:"name" validate:"required,max=40"`
	CreatedAt string `json:"created_at" validate:"required,datetime=2006-01-02 15:04:05" example:"2024-04-14 14:18:00"`
}

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

type MerchantPut struct {
	ID        uint64 `json:"-" param:"id" query:"-" validate:"required" example:"123"`
	Name      string `json:"name" validate:"required,max=40"`
	UpdatedAt string `json:"updated_at" validate:"required,datetime=2006-01-02 15:04:05" example:"2024-04-14 14:18:00"`
}

func (m *MerchantPut) ToEntity(ctx *ctx.Ctx) (res *merchantEntity.Merchant, status *resPkg.Status) {
	datetime, err := datetime.FromString(m.UpdatedAt, datetime.FormatyyyyMMddHHmmss)
	if err != nil {
		return nil, err
	}
	return &merchantEntity.Merchant{
		ID:        m.ID,
		Name:      m.Name,
		UpdatedBy: ctx.UserClaims.UserID,
		UpdatedAt: *datetime,
	}, nil
}
