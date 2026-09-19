// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"github.com/dedyf5/resik/core/organization/dto"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/organization"
	"github.com/dedyf5/resik/entities/organization/param"
	"github.com/dedyf5/resik/pkg/response"
)

//go:generate mockgen -source organization.go -package mock -destination ./mock/organization.go
type IService interface {
	OrganizationInsert(ctx *ctx.Ctx, organization *organization.Organization) (ok bool, err *response.Status)
	OrganizationUpdate(ctx *ctx.Ctx, organization *organization.Organization) (ok bool, err *response.Status)
	OrganizationGetByID(ctx *ctx.Ctx, organizationID uint64) (organization *dto.Organization, err *response.Status)
	OrganizationsGet(param *param.OrganizationsGet) (res *dto.OrganizationsResult, err *response.Status)
	OrganizationDelete(ctx *ctx.Ctx, param *organization.Organization) (ok bool, err *response.Status)
}
