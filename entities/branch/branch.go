// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package branch

import (
	"net/http"

	commonEntity "github.com/dedyf5/resik/entities/common"
	"github.com/dedyf5/resik/entities/organization"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

const TABLE_NAME = "branches"

type Branch struct {
	commonEntity.TableIdentity
	commonEntity.TableAuthor
	OrganizationID uint64                    `json:"organization_id" gorm:"column:organization_id;type:bigint(20);not null;"`
	Name           string                    `json:"name" gorm:"column:name;type:varchar(40);not null;"`
	Organization   organization.Organization `json:"organization" gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;"`
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

func (Branch) TableName() string {
	return TABLE_NAME
}
