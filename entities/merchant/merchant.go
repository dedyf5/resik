// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"time"

	"github.com/dedyf5/resik/entities/user"
)

type Merchant struct {
	ID           uint64    `json:"id"`
	UserID       uint64    `json:"user_id" gorm:"not null"`
	MerchantName string    `json:"merchant_name" gorm:"type:varchar(40);not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	CreatedBy    uint64    `json:"created_by" gorm:"not null"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamp;default:current_timestamp();not null;"`
	UpdatedBy    uint64    `json:"updated_by" gorm:"not null"`
	User         user.User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}
