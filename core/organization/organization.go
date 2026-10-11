// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"github.com/dedyf5/resik/core/organization/dto"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/organization/param"
	"github.com/dedyf5/resik/pkg/response"
)

//go:generate mockgen -source organization.go -package mock -destination ./mock/organization.go
type IService interface {
	OrganizationGetByID(ctx *ctx.Ctx, organizationID uint64) (organization *dto.Organization, err *response.Status)
	OrganizationsGet(param *param.OrganizationsGet) (res *dto.OrganizationsResult, err *response.Status)
}
