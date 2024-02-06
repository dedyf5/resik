// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"net/http"
	"time"

	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
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
			"merchant_name": "m1.merchant_name",
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
		m1.merchant_name
		`).
		Table("transaction AS t1").
		Joins("INNER JOIN merchant AS m1 ON m1.id = t1.merchant_id").
		Where("t1.merchant_id = ?", param.MerchantID).
		Where("t1.created_at >= ? AND t1.created_at <= ?", param.GroupPeriod.DatetimeStart, param.GroupPeriod.DatetimeEnd).
		Group("t1.merchant_id, period")
	if search := param.Filter.Search; search != "" {
		query = query.Where("m1.merchant_name LIKE ?", "%"+search+"%")
	}
	return
}

func (r *TransactionRepo) OutletOmzet(outletID int64, date time.Time) (*trxEntity.OutletOmzet, error) {
	return nil, nil
}

func (r *TransactionRepo) GetMerchantByID(merchantID int64) (*merchantEntity.Merchant, error) {
	return nil, nil
}

func (r *TransactionRepo) GetMerchantByIDAndUserID(merchantID int64, userID int64) (*merchantEntity.Merchant, error) {
	return nil, nil
}

func (r *TransactionRepo) GetOutletByID(outletID int64) (*outletEntity.Outlet, error) {
	return nil, nil
}

func (r *TransactionRepo) GetOutletByIDAndCreatedBy(outletID int64, createdBy int64) (*outletEntity.Outlet, error) {
	return nil, nil
}

func (r *TransactionRepo) GetUserByUserNameAndPassword(userName string, password string) (*userEntity.User, error) {
	return nil, nil
}
