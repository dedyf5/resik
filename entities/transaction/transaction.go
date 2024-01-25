// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"time"

	"github.com/dedyf5/resik/entities/merchant"
	"github.com/dedyf5/resik/entities/outlet"
)

type Transaction struct {
	ID         int64             `json:"id"`
	MerchantID int64             `json:"merchant_id" gorm:"not null"`
	OutletID   int64             `json:"outlet_id" gorm:"not null"`
	BillTotal  float64           `json:"bill_total" gorm:"not null"`
	CreatedAt  time.Time         `json:"created_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	CreatedBy  int64             `json:"created_by" gorm:"not null"`
	UpdatedAt  time.Time         `json:"updated_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	UpdatedBy  int64             `json:"updated_by" gorm:"not null"`
	Merchant   merchant.Merchant `json:"merchant" gorm:"constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
	Outlet     outlet.Outlet     `json:"outlet" gorm:"constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}
