// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"errors"
	"net/http"

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

func (r *UserRepo) MerchantIDsByUserIDGetData(userID int64) (merchantIDs []int64, status *resPkg.Status) {
	query := r.DB.Select("id").Table("merchant").Where("user_id = ?", userID)
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
