// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package req

import (
	ctx "github.com/dedyf5/resik/ctx"
	commonEntity "github.com/dedyf5/resik/entities/common"
	"github.com/dedyf5/resik/entities/groupperiod"
	trxParam "github.com/dedyf5/resik/entities/transaction/param"
	"github.com/dedyf5/resik/pkg/goku"
)

type MerchantOmzetGet struct {
	commonEntity.Request
	MerchantID    uint64  `json:"-" param:"id" query:"-" validate:"required,min=1" example:"1"`
	Mode          string  `json:"mode" query:"mode" validate:"required,oneof=day month year" example:"day"`
	DateTimeStart string  `json:"datetime_start" query:"datetime_start" validate:"required,datetime=2006-01-02 15:04:05" example:"2024-02-01 13:45:00"`
	DateTimeEnd   string  `json:"datetime_end" query:"datetime_end" validate:"required,datetime=2006-01-02 15:04:05" example:"2024-02-01 13:45:00"`
	Search        string  `json:"search" query:"search" validate:"" example:"keyword"`
	Order         *string `json:"order" query:"order" validate:"omitempty,oneof_order=period omzet merchant_name" example:"period"`
	Page          *uint   `json:"page" query:"page" validate:"omitempty,min=1" example:"1"`
	Limit         *uint   `json:"limit" query:"limit" validate:"omitempty,min=1,max=100" example:"10"`
}

func (m *MerchantOmzetGet) ToParam(c *ctx.Ctx) *trxParam.MerchantOmzetGet {
	orderStr := "period"
	if m.Order != nil {
		orderStr = *m.Order
	}
	var page int = 1
	if m.Page != nil {
		page = int(*m.Page)
	}
	var limit int = 10
	if m.Limit != nil {
		limit = int(*m.Limit)
	}
	return &trxParam.MerchantOmzetGet{
		Ctx:        c,
		MerchantID: m.MerchantID,
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
