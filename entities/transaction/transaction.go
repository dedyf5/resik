// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package transaction defines transaction entities, collections, and omzet reporting models.
package transaction

import (
	"time"

	"github.com/dedyf5/resik/entities/branch"
	commonEntity "github.com/dedyf5/resik/entities/common"
	"github.com/dedyf5/resik/entities/organization"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

// TABLE_NAME specifies the database table name for Transaction entities.
const TABLE_NAME = "transactions"

// Transaction represents a recorded transaction associated with an organization and branch.
type Transaction struct {
	commonEntity.TableIdentity
	commonEntity.TableAuthor
	OrganizationPublicID uuidPkg.UUIDV7            `json:"organization_id" gorm:"column:organization_public_id;type:uuid;not null;"`
	OrganizationName     string                    `json:"organization_name" gorm:"column:organization_name;type:varchar(40);not null;"`
	BranchPublicID       uuidPkg.UUIDV7            `json:"branch_id" gorm:"column:branch_public_id;type:uuid;not null;"`
	BranchName           string                    `json:"branch_name" gorm:"column:branch_name;type:varchar(40);not null;"`
	BillTotal            float64                   `json:"bill_total" gorm:"column:bill_total;type:double;not null;"`
	TransactedAt         time.Time                 `json:"transacted_at" gorm:"column:transacted_at;type:datetime;not null;"`
	Organization         organization.Organization `json:"organization" gorm:"foreignKey:OrganizationPublicID;references:PublicID;constraint:OnUpdate:-,OnDelete:-;"`
	Branch               branch.Branch             `json:"branch" gorm:"foreignKey:BranchPublicID;references:PublicID;constraint:OnUpdate:-,OnDelete:-;"`
}

// GetTableAuthor returns the embedded TableAuthor audit metadata for Transaction.
func (t Transaction) GetTableAuthor() commonEntity.TableAuthor {
	return t.TableAuthor
}

// TableName returns the database table name for GORM mapping.
func (Transaction) TableName() string {
	return TABLE_NAME
}

// Transactions represents a slice of Transaction entities.
type Transactions []Transaction

// GetTableAuthors extracts TableAuthor metadata from all transactions in the slice.
func (ts Transactions) GetTableAuthors() commonEntity.TableAuthors {
	return commonEntity.ExtractTableAuthors(ts)
}
