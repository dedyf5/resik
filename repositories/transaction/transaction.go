// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"net/http"

	branchEntity "github.com/dedyf5/resik/entities/branch"
	organizationEntity "github.com/dedyf5/resik/entities/organization"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	"github.com/dedyf5/resik/pkg/goku"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"gorm.io/gorm"
)

func (r *TransactionRepo) OrganizationOmzetGetData(param *paramTrx.OrganizationOmzetGet) (res []trxEntity.OrganizationOmzet, err *resPkg.Status) {
	query, err := r.OrganizationOmzetGetQuery(param)
	if err != nil {
		return res, err
	}

	query = query.
		Limit(param.Filter.LimitOrDefault()).
		Offset(param.Filter.Offset())

	if len(param.Orders) > 0 {
		orderMap := map[string]string{
			"period":            "period",
			"omzet":             "omzet",
			"organization_name": "o1.name",
		}
		order, err := goku.OrdersQueryBuilder(param.Orders, orderMap)
		if err != nil {
			return nil, resPkg.NewStatusError(http.StatusInternalServerError, err)
		}
		query = query.Order(order)
	}

	errQuery := query.Find(&res).Error
	if errQuery != nil {
		return res, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *TransactionRepo) OrganizationOmzetGetTotal(param *paramTrx.OrganizationOmzetGet) (total int64, err *resPkg.Status) {
	query, err := r.OrganizationOmzetGetQuery(param)
	if err != nil {
		return 0, err
	}
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select("COUNT(x.organization_id)").
		Table("(?) AS x", query)
	errQuery := query.Take(&total).Error
	if errQuery != nil {
		return 0, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *TransactionRepo) OrganizationOmzetGetQuery(param *paramTrx.OrganizationOmzetGet) (query *gorm.DB, err *resPkg.Status) {
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select(`
		t1.organization_id,
		DATE_FORMAT(CONVERT_TZ(t1.created_at, 'UTC', ?), ?) period,
		SUM(t1.bill_total) AS omzet,
		o1.name AS organization_name
		`, param.GroupPeriod.Timezone, param.GroupPeriod.Mode.DateFormatMySQL()).
		Table(trxEntity.TABLE_NAME+" AS t1").
		Joins("INNER JOIN "+organizationEntity.TABLE_NAME+" AS o1 ON o1.id = t1.organization_id").
		Where("t1.organization_id = ?", param.OrganizationID).
		Where("t1.created_at >= ? AND t1.created_at < ?", param.GroupPeriod.DatetimeStartString(), param.GroupPeriod.DatetimeEndString()).
		Group("t1.organization_id, period")
	if search := param.Filter.Search; search != "" {
		query = query.Where("o1.name LIKE ?", "%"+search+"%")
	}
	return
}

func (r *TransactionRepo) BranchOmzetGetData(param *paramTrx.BranchOmzetGet) (res []trxEntity.BranchOmzet, err *resPkg.Status) {
	query, err := r.BranchOmzetGetQuery(param)
	if err != nil {
		return res, err
	}

	query = query.
		Limit(param.Filter.LimitOrDefault()).
		Offset(param.Filter.Offset())

	if len(param.Orders) > 0 {
		orderMap := map[string]string{
			"period":            "period",
			"omzet":             "omzet",
			"organization_name": "o1.name",
			"branch_name":       "b1.name",
		}
		order, err := goku.OrdersQueryBuilder(param.Orders, orderMap)
		if err != nil {
			return nil, resPkg.NewStatusError(http.StatusInternalServerError, err)
		}
		query = query.Order(order)
	}

	errQuery := query.Find(&res).Error
	if errQuery != nil {
		return res, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *TransactionRepo) BranchOmzetGetTotal(param *paramTrx.BranchOmzetGet) (total int64, err *resPkg.Status) {
	query, status := r.BranchOmzetGetQuery(param)
	if status != nil {
		return 0, status
	}
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select("COUNT(x.branch_id)").
		Table("(?) AS x", query)
	errQuery := query.Take(&total).Error
	if errQuery != nil {
		return 0, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *TransactionRepo) BranchOmzetGetQuery(param *paramTrx.BranchOmzetGet) (query *gorm.DB, err *resPkg.Status) {
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select(`
		t1.organization_id,
		DATE_FORMAT(CONVERT_TZ(t1.created_at, 'UTC', ?), ?) period,
		SUM(t1.bill_total) AS omzet,
		o1.name AS organization_name,
		t1.branch_id,
		b1.name AS branch_name
		`, param.GroupPeriod.Timezone, param.GroupPeriod.Mode.DateFormatMySQL()).
		Table(trxEntity.TABLE_NAME+" AS t1").
		Joins("INNER JOIN "+organizationEntity.TABLE_NAME+" AS o1 ON o1.id = t1.organization_id").
		Joins("INNER JOIN "+branchEntity.TABLE_NAME+" AS b1 ON b1.id = t1.branch_id").
		Where("t1.branch_id = ?", param.BranchID).
		Where("t1.created_at >= ? AND t1.created_at < ?", param.GroupPeriod.DatetimeStartString(), param.GroupPeriod.DatetimeEndString()).
		Group("t1.branch_id, period")
	if search := param.Filter.Search; search != "" {
		query = query.Where("o1.name LIKE ? OR b1.name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	return
}
