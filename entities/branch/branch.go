// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package branch

import (
	"net/http"
	"time"

	"github.com/dedyf5/resik/entities/organization"
	resPkg "github.com/dedyf5/resik/pkg/response"
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
	Organization      organization.Organization `json:"organization" gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
}

func (b *Branch) BeforeCreate(tx *gorm.DB) (err error) {
	b.PublicID, err = uuidPkg.NewUUIDV7()
	return
}

func (b *Branch) IsBelongToOrganization(organizationPublicID *uuidPkg.UUIDV7) (ok bool, err *resPkg.Status) {
	if b.Organization.PublicID == *organizationPublicID {
		return true, nil
	}
	return false, resPkg.NewStatusError(
		http.StatusUnprocessableEntity,
		nil,
	)
}

type Tabler interface {
	TableName() string
}

func (Branch) TableName() string {
	return TABLE_NAME
}
