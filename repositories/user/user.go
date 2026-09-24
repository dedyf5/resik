// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"errors"
	"net/http"

	"github.com/dedyf5/resik/ctx"
	userEntity "github.com/dedyf5/resik/entities/user"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"gorm.io/gorm"
)

func (r *UserRepo) UserByID(ctx *ctx.Ctx, userID uint64) (user *userEntity.User, err *resPkg.Status) {
	errQuery := r.DB.WithContext(ctx.Context).First(&user, "id = ?", userID).Error
	if errQuery != nil {
		if errors.Is(errQuery, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

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

func (r *UserRepo) UsersGetByPublicIDs(ctx *ctx.Ctx, userIDs []uuidPkg.UUIDV7) (users userEntity.Users, err *resPkg.Status) {
	n := len(userIDs)
	if n == 0 {
		return userEntity.Users{}, nil
	}

	query := r.DB.WithContext(ctx.Context).
		Table(userEntity.TABLE_NAME)

	if len(userIDs) == 1 {
		query = query.Where("public_id = ?", userIDs[0])
	} else {
		query = query.Where("public_id IN ?", userIDs)
	}

	errQuery := query.Find(&users).Error
	if errQuery != nil {
		if errors.Is(errQuery, gorm.ErrRecordNotFound) {
			return
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}
