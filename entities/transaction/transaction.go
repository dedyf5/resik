// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"time"

	"github.com/dedyf5/resik/entities/branch"
	commonEntity "github.com/dedyf5/resik/entities/common"
	"github.com/dedyf5/resik/entities/organization"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"gorm.io/gorm"
)

const TABLE_NAME = "transactions"

type Transaction struct {
	ID                   uint64                    `json:"-" gorm:"primaryKey;autoIncrement;"`
	PublicID             uuidPkg.UUIDV7            `json:"id" gorm:"column:public_id;type:uuid;unique;not null;"`
	OrganizationPublicID uuidPkg.UUIDV7            `json:"organization_id" gorm:"column:organization_public_id;type:uuid;not null;"`
	OrganizationName     string                    `json:"organization_name" gorm:"column:organization_name;type:varchar(40);not null;"`
	BranchPublicID       uuidPkg.UUIDV7            `json:"branch_id" gorm:"column:branch_public_id;type:uuid;not null;"`
	BranchName           string                    `json:"branch_name" gorm:"column:branch_name;type:varchar(40);not null;"`
	BillTotal            float64                   `json:"bill_total" gorm:"not null;"`
	TransactedAt         time.Time                 `json:"transacted_at" gorm:"column:transacted_at;type:datetime;not null;"`
	CreatedAt            time.Time                 `json:"created_at" gorm:"type:datetime;not null;"`
	CreatedByPublicID    uuidPkg.UUIDV7            `json:"created_by" gorm:"column:created_by_public_id;type:uuid;not null;"`
	UpdatedAt            time.Time                 `json:"updated_at" gorm:"type:datetime;not null;"`
	UpdatedByPublicID    uuidPkg.UUIDV7            `json:"updated_by" gorm:"column:updated_by_public_id;type:uuid;not null;"`
	Organization         organization.Organization `json:"organization" gorm:"foreignKey:OrganizationPublicID;references:PublicID;constraint:OnUpdate:-,OnDelete:-;"`
	Branch               branch.Branch             `json:"branch" gorm:"foreignKey:BranchPublicID;references:PublicID;constraint:OnUpdate:-,OnDelete:-;"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.PublicID, err = uuidPkg.NewUUIDV7()
	return
}

func (t *Transaction) AllUserPublicIDs() []uuidPkg.UUIDV7 {
	return []uuidPkg.UUIDV7{t.CreatedByPublicID, t.UpdatedByPublicID}
}

func (t *Transaction) UniqueAllUserPublicIDs() []uuidPkg.UUIDV7 {
	return commonEntity.UniqueAllUserPublicIDs([]Transaction{*t})
}

func (Transaction) TableName() string {
	return TABLE_NAME
}

type Tabler interface {
	TableName() string
}

type Transactions []Transaction

func (ts Transactions) UniqueAllUserPublicIDs() []uuidPkg.UUIDV7 {
	return commonEntity.UniqueAllUserPublicIDs(ts)
}
