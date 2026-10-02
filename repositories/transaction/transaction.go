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

func (r *TransactionRepo) TransactionsGetData(param *paramTrx.TransactionsGet) (transactions trxEntity.Transactions, err *resPkg.Status) {
	query := r.TransactionsBaseQuery(param)

	query = query.Select("*").
		Limit(param.Filter.LimitOrDefault()).
		Offset(param.Filter.Offset())

	if len(param.Orders) > 0 {
		orderMap := map[string]string{
			"branch_name":       "branch_name",
			"organization_name": "organization_name",
			"transacted_at":     "transacted_at",
		}
		order, err := goku.OrdersQueryBuilder(param.Orders, orderMap)
		if err != nil {
			return nil, resPkg.NewStatusError(http.StatusInternalServerError, err)
		}
		query = query.Order(order)
	}

	errQuery := query.Find(&transactions).Error
	if errQuery != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *TransactionRepo) TransactionsGetTotal(param *paramTrx.TransactionsGet) (total int64, err *resPkg.Status) {
	query := r.TransactionsBaseQuery(param).Select("COUNT(id) AS total")
	errQuery := query.Take(&total).Error
	if errQuery != nil {
		return 0, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *TransactionRepo) TransactionsBaseQuery(param *paramTrx.TransactionsGet) (query *gorm.DB) {
	query = r.DB.WithContext(param.Ctx.Context).Table(trxEntity.TABLE_NAME)

	if len(param.BranchPublicIDs) == 1 {
		query = query.Where("branch_public_id = ?", param.BranchPublicIDs[0])
	} else {
		query = query.Where("branch_public_id IN ?", param.BranchPublicIDs)
	}

	if search := param.Filter.FulltextSeach(); search != "" {
		query = query.Where("MATCH(organization_name, branch_name) AGAINST(? IN BOOLEAN MODE)", search)
	}

	return
}

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
			"period": "period",
			"omzet":  "omzet",
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
		Select("COUNT(x.organization_public_id)").
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
		t1.organization_public_id,
		DATE_FORMAT(CONVERT_TZ(t1.created_at, 'UTC', ?), ?) period,
		SUM(t1.bill_total) AS omzet,
		o1.name AS organization_name
		`, param.GroupPeriod.Timezone, param.GroupPeriod.Mode.DateFormatMySQL()).
		Table(trxEntity.TABLE_NAME+" AS t1").
		Joins("LEFT JOIN "+organizationEntity.TABLE_NAME+" AS o1 ON o1.public_id = t1.organization_public_id").
		Where("t1.organization_public_id = ?", param.OrganizationPublicID).
		Where("t1.created_at >= ? AND t1.created_at < ?", param.GroupPeriod.DatetimeStartString(), param.GroupPeriod.DatetimeEndString()).
		Group("t1.organization_public_id, period")
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
			"period": "period",
			"omzet":  "omzet",
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
		Select("COUNT(x.branch_public_id)").
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
		t1.organization_public_id,
		t1.branch_public_id,
		DATE_FORMAT(CONVERT_TZ(t1.created_at, 'UTC', ?), ?) period,
		SUM(t1.bill_total) AS omzet,
		o1.name AS organization_name,
		b1.name AS branch_name
		`, param.GroupPeriod.Timezone, param.GroupPeriod.Mode.DateFormatMySQL()).
		Table(trxEntity.TABLE_NAME+" AS t1").
		Joins("LEFT JOIN "+organizationEntity.TABLE_NAME+" AS o1 ON o1.public_id = t1.organization_public_id").
		Joins("LEFT JOIN "+branchEntity.TABLE_NAME+" AS b1 ON b1.public_id = t1.branch_public_id").
		Where("t1.branch_public_id = ?", param.BranchPublicID).
		Where("t1.created_at >= ? AND t1.created_at < ?", param.GroupPeriod.DatetimeStartString(), param.GroupPeriod.DatetimeEndString()).
		Group("t1.branch_public_id, period")
	if search := param.Filter.Search; search != "" {
		query = query.Where("o1.name LIKE ? OR b1.name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	return
}
