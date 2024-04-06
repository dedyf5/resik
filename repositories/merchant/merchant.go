// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"net/http"

	"github.com/dedyf5/resik/ctx"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	paramMerchant "github.com/dedyf5/resik/entities/merchant/param"
	"github.com/dedyf5/resik/pkg/goku"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"gorm.io/gorm"
)

func (r *MerchantRepo) MerchantInsert(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, status *resPkg.Status) {
	result := r.DB.WithContext(ctx.Context).Create(merchant)
	if result.Error != nil {
		return false, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: result.Error,
		}
	}
	return true, nil
}

func (r *MerchantRepo) MerchantUpdate(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, status *resPkg.Status) {
	result := r.DB.WithContext(ctx.Context).
		Exec("UPDATE "+merchantEntity.TABLE_NAME+" SET name = ?, updated_at = ?, updated_by = ? WHERE id = ?", merchant.Name, merchant.UpdatedAt, merchant.UpdatedBy, merchant.ID)
	if result.Error != nil {
		return false, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: result.Error,
		}
	}
	return true, nil
}

func (r *MerchantRepo) MerchantListGetData(param *paramMerchant.MerchantListGet) (merchant []merchantEntity.Merchant, status *resPkg.Status) {
	query := r.MerchantListBaseQuery(param)

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
			return nil, &resPkg.Status{
				Code:       http.StatusInternalServerError,
				CauseError: err,
			}
		}
		query = query.Order(order)
	}

	err := query.Find(&merchant).Error
	if err != nil {
		return nil, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	return
}

func (r *MerchantRepo) MerchantListGetTotal(param *paramMerchant.MerchantListGet) (total uint64, status *resPkg.Status) {
	query := r.MerchantListBaseQuery(param).Select("COUNT(id) AS total")
	err := query.Take(&total).Error
	if err != nil {
		return 0, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	return
}

func (r *MerchantRepo) MerchantListBaseQuery(param *paramMerchant.MerchantListGet) (query *gorm.DB) {
	query = r.DB.WithContext(param.Ctx.Context).
		Table(merchantEntity.TABLE_NAME).
		Where("id IN ?", param.MerchantIDs)
	if param.Filter.Search != "" {
		query = query.Where("name LIKE ?", "%"+param.Filter.Search+"%")
	}
	return
}
