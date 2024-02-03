// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	ctx "github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/groupperiod"
	trxParam "github.com/dedyf5/resik/entities/transaction/param"
	"github.com/dedyf5/resik/pkg/goku"
)

type MerchantOmzetGet struct {
	ID            uint64  `json:"id" param:"id" validate:"required,min=1" example:"1"`
	Mode          string  `json:"mode" query:"mode" validate:"required,oneof=day month year" example:"day"`
	DateTimeStart string  `json:"datetime_start" query:"datetime_start" validate:"required,datetime=2006-01-02 15:04:05" example:"2024-02-01 13:45:00"`
	DateTimeEnd   string  `json:"datetime_end" query:"datetime_end" validate:"required,datetime=2006-01-02 15:04:05" example:"2024-02-01 13:45:00"`
	Search        string  `json:"search" query:"search" validate:"" example:"keyword"`
	Order         *string `json:"order" query:"order" validate:"omitempty" example:"period"`
	Page          *uint   `json:"page" query:"page" validate:"omitempty,min=1" example:"1"`
	Limit         *uint   `json:"limit" query:"limit" validate:"omitempty,min=1,max=100" example:"10"`
}

type TransactionUpsert struct {
	Name      string `json:"name" validate:"required" example:"Shiina"`
	Address   string `json:"address" validate:"required" example:"Jl. Nusantara"`
	QTY       int64  `json:"qty" validate:"required,min=1" example:"3"`
	CreatedAt string `json:"created_at" validate:"required,datetime=2006-01-02 15:04:05" example:"2024-02-01 13:45:00"`
}

func (m *MerchantOmzetGet) ToParam(c *ctx.Ctx) *trxParam.MerchantOmzetGet {
	orderStr := "period"
	if m.Order != nil {
		orderStr = *m.Order
	}
	var page uint = 1
	if m.Page != nil {
		page = *m.Page
	}
	var limit uint = 10
	if m.Limit != nil {
		limit = *m.Limit
	}
	return &trxParam.MerchantOmzetGet{
		Ctx:        c,
		MerchantID: m.ID,
		GroupPeriod: groupperiod.GroupPeriod{
			Mode:          groupperiod.Mode(m.Mode),
			DatetimeStart: m.DateTimeStart,
			DatetimeEnd:   m.DateTimeEnd,
		},
		Filter: goku.Filter{
			Search: m.Search,
			Page:   page,
			Limit:  limit,
		},
		Orders: goku.OrdersBuilder(orderStr),
	}
}
