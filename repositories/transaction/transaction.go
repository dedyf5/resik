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
	statusPkg "github.com/dedyf5/resik/pkg/status"
)

func (r *TransactionRepo) MerchantOmzetGet(param *paramTrx.MerchantOmzetGet) (res []trxEntity.MerchantOmzet, status *statusPkg.Status) {
	query := r.DB.
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
		Group("t1.merchant_id, period").
		Limit(param.Filter.LimitOrDefault()).
		Offset(param.Filter.Offset())
	if search := param.Filter.Search; search != "" {
		query = query.Where("m1.merchant_name LIKE ?", "%"+search+"%")
	}
	if len(param.Orders) > 0 {
		orderMap := map[string]string{
			"period":        "period",
			"omzet":         "omzet",
			"merchant_name": "m1.merchant_name",
		}
		order, err := goku.OrdersQueryBuilder(param.Orders, orderMap)
		if err != nil {
			return res, &statusPkg.Status{
				Code:       http.StatusInternalServerError,
				CauseError: err,
			}
		}
		query = query.Order(order)
	}
	err := query.Find(&res).Error
	if err != nil {
		return res, &statusPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
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
