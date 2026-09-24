// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package repositories

import (
	"github.com/dedyf5/resik/ctx"
	checkEntity "github.com/dedyf5/resik/entities/check"
	organizationEntity "github.com/dedyf5/resik/entities/organization"
	paramOrganization "github.com/dedyf5/resik/entities/organization/param"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

//go:generate mockgen -source repository.go -package mock -destination ./mock/repository.go
type ICheck interface {
	Check() checkEntity.CheckDetail
}

type ITransaction interface {
	OrganizationOmzetGetData(param *paramTrx.OrganizationOmzetGet) (res []trxEntity.OrganizationOmzet, err *resPkg.Status)
	OrganizationOmzetGetTotal(param *paramTrx.OrganizationOmzetGet) (total int64, err *resPkg.Status)
	BranchOmzetGetData(param *paramTrx.BranchOmzetGet) (res []trxEntity.BranchOmzet, err *resPkg.Status)
	BranchOmzetGetTotal(param *paramTrx.BranchOmzetGet) (total int64, err *resPkg.Status)
}

type IUser interface {
	UserByID(ctx *ctx.Ctx, userID uint64) (user *userEntity.User, err *resPkg.Status)
	UserByUsername(ctx *ctx.Ctx, username string) (user *userEntity.User, err *resPkg.Status)
	UsersGetByPublicIDs(ctx *ctx.Ctx, userIDs []uuidPkg.UUIDV7) (users userEntity.Users, err *resPkg.Status)
}

type IOrganization interface {
	OrganizationInsert(ctx *ctx.Ctx, organization *organizationEntity.Organization) (ok bool, err *resPkg.Status)
	OrganizationUpdate(ctx *ctx.Ctx, organization *organizationEntity.Organization) (ok bool, err *resPkg.Status)
	OrganizationGetByID(ctx *ctx.Ctx, organizationID uint64) (organization *organizationEntity.Organization, err *resPkg.Status)
	OrganizationsGetData(param *paramOrganization.OrganizationsGet) (organizations organizationEntity.Organizations, err *resPkg.Status)
	OrganizationsGetTotal(param *paramOrganization.OrganizationsGet) (total int64, err *resPkg.Status)
	OrganizationDelete(c *ctx.Ctx, organization *organizationEntity.Organization) (ok bool, err *resPkg.Status)
}
