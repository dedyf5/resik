// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package outlet

import (
	"time"

	"github.com/dedyf5/resik/entities/merchant"
)

type Outlet struct {
	ID         int64             `json:"id"`
	MerchantID int64             `json:"merchant_id" gorm:"not null"`
	OutletName string            `json:"outlet_name" gorm:"type:varchar(40);not null"`
	CreatedAt  time.Time         `json:"created_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	CreatedBy  int64             `json:"created_by" gorm:"not null"`
	UpdatedAt  time.Time         `json:"updated_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	UpdatedBy  int64             `json:"updated_by" gorm:"not null"`
	Merchant   merchant.Merchant `json:"merchant" gorm:"constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}
