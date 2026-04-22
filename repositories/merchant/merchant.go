// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"errors"
	"net/http"

	"github.com/dedyf5/resik/ctx"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	paramMerchant "github.com/dedyf5/resik/entities/merchant/param"
	"github.com/dedyf5/resik/pkg/goku"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"gorm.io/gorm"
)

func (r *MerchantRepo) MerchantInsert(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, err *resPkg.Status) {
	result := r.DB.WithContext(ctx.Context).Create(merchant)
	if result.Error != nil {
		return false, resPkg.NewStatusError(http.StatusInternalServerError, result.Error)
	}
	return true, nil
}

func (r *MerchantRepo) MerchantUpdate(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, err *resPkg.Status) {
	result := r.DB.WithContext(ctx.Context).
		Exec("UPDATE "+merchantEntity.TABLE_NAME+" SET name = ?, description = ?, updated_at = ?, updated_by = ? WHERE id = ?", merchant.Name, merchant.Description, merchant.UpdatedAt, merchant.UpdatedBy, merchant.ID)
	if result.Error != nil {
		return false, resPkg.NewStatusError(http.StatusInternalServerError, result.Error)
	}
	return true, nil
}

func (r *MerchantRepo) MerchantGetByIDAndUserID(ctx *ctx.Ctx, merchantID, userID uint64) (merchant *merchantEntity.Merchant, err *resPkg.Status) {
	errDB := r.DB.WithContext(ctx.Context).Where("id = ? AND user_id = ?", merchantID, userID).First(&merchant).Error
	if errDB != nil {
		if errors.Is(errDB, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errDB)
	}
	return
}

func (r *MerchantRepo) MerchantListGetData(param *paramMerchant.MerchantListGet) (merchant []merchantEntity.Merchant, err *resPkg.Status) {
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
			return nil, resPkg.NewStatusError(http.StatusInternalServerError, err)
		}
		query = query.Order(order)
	}

	errQuery := query.Find(&merchant).Error
	if errQuery != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *MerchantRepo) MerchantListGetTotal(param *paramMerchant.MerchantListGet) (total int64, err *resPkg.Status) {
	query := r.MerchantListBaseQuery(param).Select("COUNT(id) AS total")
	errQuery := query.Take(&total).Error
	if errQuery != nil {
		return 0, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
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

func (r *MerchantRepo) MerchantDelete(c *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, err *resPkg.Status) {
	if err := r.DB.WithContext(c.Context).Delete(merchant).Error; err != nil {
		return false, resPkg.NewStatusError(http.StatusInternalServerError, err)
	}
	return true, nil
}
