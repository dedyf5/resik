// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"net/http"

	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	"github.com/dedyf5/resik/pkg/goku"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"gorm.io/gorm"
)

func (r *TransactionRepo) MerchantOmzetGetData(param *paramTrx.MerchantOmzetGet) (res []trxEntity.MerchantOmzet, err *resPkg.Status) {
	query, err := r.MerchantOmzetGetQuery(param)
	if err != nil {
		return res, err
	}

	query = query.
		Limit(param.Filter.LimitOrDefault()).
		Offset(param.Filter.Offset())

	if len(param.Orders) > 0 {
		orderMap := map[string]string{
			"period":        "period",
			"omzet":         "omzet",
			"merchant_name": "m1.name",
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

func (r *TransactionRepo) MerchantOmzetGetTotal(param *paramTrx.MerchantOmzetGet) (total int64, err *resPkg.Status) {
	query, err := r.MerchantOmzetGetQuery(param)
	if err != nil {
		return 0, err
	}
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select("COUNT(x.merchant_id)").
		Table("(?) AS x", query)
	errQuery := query.Take(&total).Error
	if errQuery != nil {
		return 0, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *TransactionRepo) MerchantOmzetGetQuery(param *paramTrx.MerchantOmzetGet) (query *gorm.DB, err *resPkg.Status) {
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select(`
		t1.merchant_id,
		DATE_FORMAT(CONVERT_TZ(t1.created_at, 'UTC', ?), ?) period,
		SUM(t1.bill_total) AS omzet,
		m1.name AS merchant_name
		`, param.GroupPeriod.Timezone, param.GroupPeriod.Mode.DateFormatMySQL()).
		Table(trxEntity.TABLE_NAME+" AS t1").
		Joins("INNER JOIN "+merchantEntity.TABLE_NAME+" AS m1 ON m1.id = t1.merchant_id").
		Where("t1.merchant_id = ?", param.MerchantID).
		Where("t1.created_at >= ? AND t1.created_at < ?", param.GroupPeriod.DatetimeStartString(), param.GroupPeriod.DatetimeEndString()).
		Group("t1.merchant_id, period")
	if search := param.Filter.Search; search != "" {
		query = query.Where("m1.name LIKE ?", "%"+search+"%")
	}
	return
}

func (r *TransactionRepo) OutletOmzetGetData(param *paramTrx.OutletOmzetGet) (res []trxEntity.OutletOmzet, err *resPkg.Status) {
	query, err := r.OutletOmzetGetQuery(param)
	if err != nil {
		return res, err
	}

	query = query.
		Limit(param.Filter.LimitOrDefault()).
		Offset(param.Filter.Offset())

	if len(param.Orders) > 0 {
		orderMap := map[string]string{
			"period":        "period",
			"omzet":         "omzet",
			"merchant_name": "m1.name",
			"outlet_name":   "o1.name",
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

func (r *TransactionRepo) OutletOmzetGetTotal(param *paramTrx.OutletOmzetGet) (total int64, err *resPkg.Status) {
	query, status := r.OutletOmzetGetQuery(param)
	if status != nil {
		return 0, status
	}
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select("COUNT(x.outlet_id)").
		Table("(?) AS x", query)
	errQuery := query.Take(&total).Error
	if errQuery != nil {
		return 0, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *TransactionRepo) OutletOmzetGetQuery(param *paramTrx.OutletOmzetGet) (query *gorm.DB, err *resPkg.Status) {
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select(`
		t1.merchant_id,
		DATE_FORMAT(CONVERT_TZ(t1.created_at, 'UTC', ?), ?) period,
		SUM(t1.bill_total) AS omzet,
		m1.name AS merchant_name,
		t1.outlet_id,
		o1.name AS outlet_name
		`, param.GroupPeriod.Timezone, param.GroupPeriod.Mode.DateFormatMySQL()).
		Table(trxEntity.TABLE_NAME+" AS t1").
		Joins("INNER JOIN "+merchantEntity.TABLE_NAME+" AS m1 ON m1.id = t1.merchant_id").
		Joins("INNER JOIN "+outletEntity.TABLE_NAME+" AS o1 ON o1.id = t1.outlet_id").
		Where("t1.outlet_id = ?", param.OutletID).
		Where("t1.created_at >= ? AND t1.created_at < ?", param.GroupPeriod.DatetimeStartString(), param.GroupPeriod.DatetimeEndString()).
		Group("t1.outlet_id, period")
	if search := param.Filter.Search; search != "" {
		query = query.Where("m1.name LIKE ? OR o1.name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	return
}

func (r *TransactionRepo) GetMerchantByID(merchantID uint64) (*merchantEntity.Merchant, error) {
	return nil, nil
}

func (r *TransactionRepo) GetMerchantByIDAndUserID(merchantID uint64, userID uint64) (*merchantEntity.Merchant, error) {
	return nil, nil
}

func (r *TransactionRepo) GetOutletByID(outletID uint64) (*outletEntity.Outlet, error) {
	return nil, nil
}

func (r *TransactionRepo) GetOutletByIDAndCreatedBy(outletID uint64, createdBy uint64) (*outletEntity.Outlet, error) {
	return nil, nil
}
