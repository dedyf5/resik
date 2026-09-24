// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	"errors"
	"net/http"

	"github.com/dedyf5/resik/ctx"
	organizationEntity "github.com/dedyf5/resik/entities/organization"
	paramOrganization "github.com/dedyf5/resik/entities/organization/param"
	"github.com/dedyf5/resik/pkg/goku"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"gorm.io/gorm"
)

func (r *OrganizationRepo) OrganizationInsert(ctx *ctx.Ctx, organization *organizationEntity.Organization) (ok bool, err *resPkg.Status) {
	result := r.DB.WithContext(ctx.Context).Create(organization)
	if result.Error != nil {
		return false, resPkg.NewStatusError(http.StatusInternalServerError, result.Error)
	}
	return true, nil
}

func (r *OrganizationRepo) OrganizationUpdate(ctx *ctx.Ctx, organization *organizationEntity.Organization) (ok bool, err *resPkg.Status) {
	result := r.DB.WithContext(ctx.Context).
		Exec("UPDATE "+organizationEntity.TABLE_NAME+" SET name = ?, description = ?, updated_at = ?, updated_by_public_id = ? WHERE id = ?", organization.Name, organization.Description, organization.UpdatedAt, organization.UpdatedByPublicID, organization.ID)
	if result.Error != nil {
		return false, resPkg.NewStatusError(http.StatusInternalServerError, result.Error)
	}
	return true, nil
}

func (r *OrganizationRepo) OrganizationGetByID(ctx *ctx.Ctx, organizationID uint64) (organization *organizationEntity.Organization, err *resPkg.Status) {
	errDB := r.DB.WithContext(ctx.Context).Where("id = ?", organizationID).First(&organization).Error
	if errDB != nil {
		if errors.Is(errDB, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errDB)
	}
	return
}

func (r *OrganizationRepo) OrganizationsGetData(param *paramOrganization.OrganizationsGet) (organizations organizationEntity.Organizations, err *resPkg.Status) {
	query := r.OrganizationsBaseQuery(param)

	query = query.Select("*").
		Limit(param.Filter.LimitOrDefault()).
		Offset(param.Filter.Offset())

	if len(param.Orders) > 0 {
		orderMap := map[string]string{
			"name":       "name",
			"created_at": "created_at",
			"updated_at": "updated_at",
		}
		order, err := goku.OrdersQueryBuilder(param.Orders, orderMap)
		if err != nil {
			return nil, resPkg.NewStatusError(http.StatusInternalServerError, err)
		}
		query = query.Order(order)
	}

	errQuery := query.Find(&organizations).Error
	if errQuery != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *OrganizationRepo) OrganizationsGetTotal(param *paramOrganization.OrganizationsGet) (total int64, err *resPkg.Status) {
	query := r.OrganizationsBaseQuery(param).Select("COUNT(id) AS total")
	errQuery := query.Take(&total).Error
	if errQuery != nil {
		return 0, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *OrganizationRepo) OrganizationsBaseQuery(param *paramOrganization.OrganizationsGet) (query *gorm.DB) {
	query = r.DB.WithContext(param.Ctx.Context).
		Table(organizationEntity.TABLE_NAME)

	if len(param.OrganizationIDs) == 1 {
		query = query.Where("id = ?", param.OrganizationIDs[0])
	} else {
		query = query.Where("id IN ?", param.OrganizationIDs)
	}

	if param.Filter.Search != "" {
		query = query.Where("name LIKE ?", "%"+param.Filter.Search+"%")
	}
	return
}

func (r *OrganizationRepo) OrganizationDelete(c *ctx.Ctx, organization *organizationEntity.Organization) (ok bool, err *resPkg.Status) {
	if err := r.DB.WithContext(c.Context).Delete(organization).Error; err != nil {
		return false, resPkg.NewStatusError(http.StatusInternalServerError, err)
	}
	return true, nil
}
