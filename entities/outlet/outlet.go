// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package outlet

import (
	"time"

	"github.com/dedyf5/resik/entities/merchant"
)

const TABLE_NAME = "outlet"

type Outlet struct {
	ID         uint64            `json:"id"`
	MerchantID uint64            `json:"merchant_id" gorm:"not null"`
	Name       string            `json:"name" gorm:"type:varchar(40);not null"`
	CreatedAt  time.Time         `json:"created_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	CreatedBy  uint64            `json:"created_by" gorm:"not null"`
	UpdatedAt  time.Time         `json:"updated_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	UpdatedBy  uint64            `json:"updated_by" gorm:"not null"`
	Merchant   merchant.Merchant `json:"merchant" gorm:"constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}

type Tabler interface {
	TableName() string
}

func (Outlet) TableName() string {
	return TABLE_NAME
}
