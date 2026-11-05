// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/drivers"
	branchEntity "github.com/dedyf5/resik/entities/branch"
	organizationEntity "github.com/dedyf5/resik/entities/organization"
	paramOrganization "github.com/dedyf5/resik/entities/organization/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"github.com/dedyf5/resik/repositories/organization/mysql"
	"github.com/dedyf5/resik/repositories/organization/postgres"
	"gorm.io/gorm"
)

//go:generate mockgen -source organization.go -package mock -destination ./mock/organization.go
type IOrganization interface {
	OrganizationGetByID(ctx *ctx.Ctx, organizationID uint64) (organization *organizationEntity.Organization, err *resPkg.Status)
	OrganizationsGetData(param *paramOrganization.OrganizationsGet) (organizations organizationEntity.Organizations, err *resPkg.Status)
	OrganizationsGetTotal(param *paramOrganization.OrganizationsGet) (total int64, err *resPkg.Status)

	BranchGetByPublicID(ctx *ctx.Ctx, branchPublicID *uuidPkg.UUIDV7) (branch *branchEntity.Branch, err *resPkg.Status)
}

func New(db *gorm.DB, engine drivers.SQLEngine) IOrganization {
	switch engine {
	case drivers.MySQL:
		return mysql.New(db)
	case drivers.PostgreSQL:
		return postgres.New(db)
	}
	return nil
}
