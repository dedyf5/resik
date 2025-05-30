// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"errors"
	"net/http"

	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	userEntity "github.com/dedyf5/resik/entities/user"
	paramUser "github.com/dedyf5/resik/entities/user/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"gorm.io/gorm"
)

func (r *UserRepo) UserByUsernameAndPasswordGetData(param paramUser.Auth) (user *userEntity.User, status *resPkg.Status) {
	var res userEntity.User
	err := r.DB.WithContext(param.Ctx.Context).First(&res, "username = ? AND password = ?", param.Username, param.Password).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	return &res, nil
}

func (r *UserRepo) MerchantIDsByUserIDGetData(userID uint64) (merchantIDs []uint64, status *resPkg.Status) {
	query := r.DB.Select("id").Table(merchantEntity.TABLE_NAME).Where("user_id = ?", userID)
	err := query.Find(&merchantIDs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		return nil, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	return
}

func (r *UserRepo) OutletMerchantByUserIDGetData(userID uint64) (outlets []outletEntity.Outlet, status *resPkg.Status) {
	query := r.DB.Select("o1.id, m1.id AS merchant_id").
		Table(outletEntity.TABLE_NAME+" AS o1").
		Joins("RIGHT JOIN "+merchantEntity.TABLE_NAME+" AS m1 ON m1.id = o1.merchant_id").
		Where("m1.user_id = ?", userID)
	err := query.Find(&outlets).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		return nil, &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}
	return
}
