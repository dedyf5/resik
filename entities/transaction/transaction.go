// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"time"

	"github.com/dedyf5/resik/entities/merchant"
	"github.com/dedyf5/resik/entities/outlet"
)

const TABLE_NAME = "transactions"

type Transaction struct {
	ID         uint64            `json:"id"`
	MerchantID uint64            `json:"merchant_id" gorm:"not null"`
	OutletID   uint64            `json:"outlet_id" gorm:"not null"`
	BillTotal  float64           `json:"bill_total" gorm:"not null"`
	CreatedAt  time.Time         `json:"created_at" gorm:"type:datetime;not null;"`
	CreatedBy  uint64            `json:"created_by" gorm:"not null"`
	UpdatedAt  time.Time         `json:"updated_at" gorm:"type:datetime;not null;"`
	UpdatedBy  uint64            `json:"updated_by" gorm:"not null"`
	Merchant   merchant.Merchant `json:"merchant" gorm:"constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
	Outlet     outlet.Outlet     `json:"outlet" gorm:"constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}

type Tabler interface {
	TableName() string
}

func (Transaction) TableName() string {
	return TABLE_NAME
}
