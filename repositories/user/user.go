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

func (r *UserRepo) UsersGetByIDs(ctx *ctx.Ctx, userIDs []uint64) (users userEntity.Users, err *resPkg.Status) {
	n := len(userIDs)
	if n == 0 {
		return userEntity.Users{}, nil
	}

	query := r.DB.WithContext(ctx.Context).
		Table(userEntity.TABLE_NAME)

	if len(userIDs) == 1 {
		query = query.Where("id = ?", userIDs[0])
	} else {
		query = query.Where("id IN ?", userIDs)
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

func (r *UserRepo) MerchantIDsByUserIDGetData(userID uint64) (merchantIDs []uint64, err *resPkg.Status) {
	query := r.DB.Select("id").Table(merchantEntity.TABLE_NAME).Where("owner_id = ?", userID)
	errQuery := query.Find(&merchantIDs).Error
	if errQuery != nil {
		if errors.Is(errQuery, gorm.ErrRecordNotFound) {
			return
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *UserRepo) OutletMerchantByUserIDGetData(ctx *ctx.Ctx, userID uint64) (merchantOutletIDs userEntity.MerchantOutletIDs, err *resPkg.Status) {
	query := r.DB.WithContext(ctx.Context).
		Select("o1.id AS outlet_id, o1.public_id AS outlet_public_id, m1.id AS merchant_id, m1.public_id AS merchant_public_id").
		Table(outletEntity.TABLE_NAME+" AS o1").
		Joins("RIGHT JOIN "+merchantEntity.TABLE_NAME+" AS m1 ON m1.id = o1.merchant_id").
		Where("m1.owner_id = ?", userID)
	errQuery := query.Find(&merchantOutletIDs).Error
	if errQuery != nil {
		if errors.Is(errQuery, gorm.ErrRecordNotFound) {
			return
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}
