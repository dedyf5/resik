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
)

const TABLE_NAME = "transactions"

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

func (t Transaction) GetTableAuthor() commonEntity.TableAuthor {
	return t.TableAuthor
}

func (Transaction) TableName() string {
	return TABLE_NAME
}

type Transactions []Transaction

func (ts Transactions) GetTableAuthors() commonEntity.TableAuthors {
	return commonEntity.ExtractTableAuthors(ts)
}
