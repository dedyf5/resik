// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"errors"
	"net/http"

	"github.com/dedyf5/resik/ctx"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	userEntity "github.com/dedyf5/resik/entities/user"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"gorm.io/gorm"
)

func (r *UserRepo) UserByUsername(ctx *ctx.Ctx, username string) (user *userEntity.User, err *resPkg.Status) {
	var res userEntity.User
	errQuery := r.DB.WithContext(ctx.Context).First(&res, "username = ?", username).Error
	if errQuery != nil {
		if errors.Is(errQuery, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return &res, nil
}

func (r *UserRepo) MerchantIDsByUserIDGetData(userID uint64) (merchantIDs []uint64, err *resPkg.Status) {
	query := r.DB.Select("id").Table(merchantEntity.TABLE_NAME).Where("user_id = ?", userID)
	errQuery := query.Find(&merchantIDs).Error
	if errQuery != nil {
		if errors.Is(errQuery, gorm.ErrRecordNotFound) {
			return
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *UserRepo) OutletMerchantByUserIDGetData(userID uint64) (outlets []outletEntity.Outlet, err *resPkg.Status) {
	query := r.DB.Select("o1.id, m1.id AS merchant_id").
		Table(outletEntity.TABLE_NAME+" AS o1").
		Joins("RIGHT JOIN "+merchantEntity.TABLE_NAME+" AS m1 ON m1.id = o1.merchant_id").
		Where("m1.user_id = ?", userID)
	errQuery := query.Find(&outlets).Error
	if errQuery != nil {
		if errors.Is(errQuery, gorm.ErrRecordNotFound) {
			return
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}
