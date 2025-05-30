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

func (r *TransactionRepo) MerchantOmzetGetData(param *paramTrx.MerchantOmzetGet) (res []trxEntity.MerchantOmzet, status *resPkg.Status) {
	query, status := r.MerchantOmzetGetQuery(param)
	if status != nil {
		return res, status
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
			return nil, &resPkg.Status{
				Code:       http.StatusInternalServerError,
				CauseError: err,
			}
		}
		query = query.Order(order)
	}

	err := query.Find(&res).Error
	if err != nil {
		return res, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	return
}

func (r *TransactionRepo) MerchantOmzetGetTotal(param *paramTrx.MerchantOmzetGet) (total uint64, status *resPkg.Status) {
	query, status := r.MerchantOmzetGetQuery(param)
	if status != nil {
		return 0, status
	}
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select("COUNT(x.merchant_id)").
		Table("(?) AS x", query)
	err := query.Take(&total).Error
	if err != nil {
		return 0, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	return
}

func (r *TransactionRepo) MerchantOmzetGetQuery(param *paramTrx.MerchantOmzetGet) (query *gorm.DB, status *resPkg.Status) {
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select(`
		t1.merchant_id,
		DATE_FORMAT(t1.created_at, '`+param.GroupPeriod.Mode.DateFormatMySQL()+`') period,
		SUM(t1.bill_total) AS omzet,
		m1.name AS merchant_name
		`).
		Table(trxEntity.TABLE_NAME+" AS t1").
		Joins("INNER JOIN "+merchantEntity.TABLE_NAME+" AS m1 ON m1.id = t1.merchant_id").
		Where("t1.merchant_id = ?", param.MerchantID).
		Where("t1.created_at >= ? AND t1.created_at <= ?", param.GroupPeriod.DatetimeStart, param.GroupPeriod.DatetimeEnd).
		Group("t1.merchant_id, period")
	if search := param.Filter.Search; search != "" {
		query = query.Where("m1.name LIKE ?", "%"+search+"%")
	}
	return
}

func (r *TransactionRepo) OutletOmzetGetData(param *paramTrx.OutletOmzetGet) (res []trxEntity.OutletOmzet, status *resPkg.Status) {
	query, status := r.OutletOmzetGetQuery(param)
	if status != nil {
		return res, status
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
			return nil, &resPkg.Status{
				Code:       http.StatusInternalServerError,
				CauseError: err,
			}
		}
		query = query.Order(order)
	}

	err := query.Find(&res).Error
	if err != nil {
		return res, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	return
}

func (r *TransactionRepo) OutletOmzetGetTotal(param *paramTrx.OutletOmzetGet) (total uint64, status *resPkg.Status) {
	query, status := r.OutletOmzetGetQuery(param)
	if status != nil {
		return 0, status
	}
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select("COUNT(x.outlet_id)").
		Table("(?) AS x", query)
	err := query.Take(&total).Error
	if err != nil {
		return 0, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	return
}

func (r *TransactionRepo) OutletOmzetGetQuery(param *paramTrx.OutletOmzetGet) (query *gorm.DB, status *resPkg.Status) {
	query = r.DB.
		WithContext(param.Ctx.Context).
		Select(`
		t1.merchant_id,
		DATE_FORMAT(t1.created_at, '`+param.GroupPeriod.Mode.DateFormatMySQL()+`') period,
		SUM(t1.bill_total) AS omzet,
		m1.name AS merchant_name,
		t1.outlet_id,
		o1.name AS outlet_name
		`).
		Table(trxEntity.TABLE_NAME+" AS t1").
		Joins("INNER JOIN "+merchantEntity.TABLE_NAME+" AS m1 ON m1.id = t1.merchant_id").
		Joins("INNER JOIN "+outletEntity.TABLE_NAME+" AS o1 ON o1.id = t1.outlet_id").
		Where("t1.outlet_id = ?", param.OutletID).
		Where("t1.created_at >= ? AND t1.created_at <= ?", param.GroupPeriod.DatetimeStart, param.GroupPeriod.DatetimeEnd).
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
