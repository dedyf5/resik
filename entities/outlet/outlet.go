// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package outlet

import (
	"time"

	"github.com/dedyf5/resik/entities/merchant"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"gorm.io/gorm"
)

const TABLE_NAME = "outlets"

type Outlet struct {
	ID         uint64            `json:"-" gorm:"primaryKey;autoIncrement;"`
	PublicID   uuidPkg.UUIDV7    `json:"id" gorm:"column:public_id;type:uuid;unique;not null;"`
	MerchantID uint64            `json:"merchant_id" gorm:"not null"`
	Name       string            `json:"name" gorm:"type:varchar(40);not null"`
	CreatedAt  time.Time         `json:"created_at" gorm:"type:datetime;not null;"`
	CreatedBy  uint64            `json:"created_by" gorm:"not null"`
	UpdatedAt  time.Time         `json:"updated_at" gorm:"type:datetime;not null;"`
	UpdatedBy  uint64            `json:"updated_by" gorm:"not null"`
	Merchant   merchant.Merchant `json:"merchant" gorm:"constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}

func (o *Outlet) BeforeCreate(tx *gorm.DB) (err error) {
	o.PublicID, err = uuidPkg.NewUUIDV7()
	return
}

type Tabler interface {
	TableName() string
}

func (Outlet) TableName() string {
	return TABLE_NAME
}
