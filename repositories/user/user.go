// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"errors"
	"net/http"

	userEntity "github.com/dedyf5/resik/entities/user"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"gorm.io/gorm"
)

func (r *UserRepo) GetUserByUsernameAndPassword(userName string, password string) (user *userEntity.User, status *resPkg.Status) {
	var res userEntity.User
	err := r.DB.First(&res, "username = ? AND password = ?", userName, password).Error
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
