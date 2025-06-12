// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package outlet

import (
	"slices"
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

func GetUniqueMerchantIDsAndOutletIDs(outlets []Outlet) (merchantIDs, OutletIDs []uint64) {
	length := len(outlets)
	MIDs := make([]uint64, 0, length)
	OIDs := make([]uint64, 0, length)
	for _, v := range outlets {
		if v.ID > 0 {
			OIDs = append(OIDs, v.ID)
		}
		if !slices.Contains(MIDs, v.MerchantID) {
			MIDs = append(MIDs, v.MerchantID)
		}
	}
	return MIDs, OIDs
}
