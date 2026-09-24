// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package branch

import (
	"time"

	"github.com/dedyf5/resik/entities/organization"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"gorm.io/gorm"
)

const TABLE_NAME = "branches"

type Branch struct {
	ID                uint64                    `json:"-" gorm:"primaryKey;autoIncrement;"`
	PublicID          uuidPkg.UUIDV7            `json:"id" gorm:"column:public_id;type:uuid;unique;not null;"`
	OrganizationID    uint64                    `json:"organization_id" gorm:"not null"`
	Name              string                    `json:"name" gorm:"type:varchar(40);not null"`
	CreatedAt         time.Time                 `json:"created_at" gorm:"type:datetime;not null;"`
	CreatedByPublicID uuidPkg.UUIDV7            `json:"created_by" gorm:"column:created_by_public_id;type:uuid;not null;"`
	UpdatedAt         time.Time                 `json:"updated_at" gorm:"type:datetime;not null;"`
	UpdatedByPublicID uuidPkg.UUIDV7            `json:"updated_by" gorm:"column:updated_by_public_id;type:uuid;not null;"`
	Organization      organization.Organization `json:"organization" gorm:"constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}

func (o *Branch) BeforeCreate(tx *gorm.DB) (err error) {
	o.PublicID, err = uuidPkg.NewUUIDV7()
	return
}

type Tabler interface {
	TableName() string
}

func (Branch) TableName() string {
	return TABLE_NAME
}
