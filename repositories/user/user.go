// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"errors"
	"net/http"

	"github.com/dedyf5/resik/ctx"
	branchEntity "github.com/dedyf5/resik/entities/branch"
	organizationEntity "github.com/dedyf5/resik/entities/organization"
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

func (r *UserRepo) OrganizationIDsByUserIDGetData(userID uint64) (organizationIDs []uint64, err *resPkg.Status) {
	query := r.DB.Select("id").Table(organizationEntity.TABLE_NAME).Where("owner_id = ?", userID)
	errQuery := query.Find(&organizationIDs).Error
	if errQuery != nil {
		if errors.Is(errQuery, gorm.ErrRecordNotFound) {
			return
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}

func (r *UserRepo) BranchOrganizationByUserIDGetData(ctx *ctx.Ctx, userID uint64) (organizationBranchIDs userEntity.OrganizationBranchIDs, err *resPkg.Status) {
	query := r.DB.WithContext(ctx.Context).
		Select("b1.id AS branch_id, b1.public_id AS branch_public_id, o1.id AS organization_id, o1.public_id AS organization_public_id").
		Table(branchEntity.TABLE_NAME+" AS b1").
		Joins("RIGHT JOIN "+organizationEntity.TABLE_NAME+" AS o1 ON o1.id = b1.organization_id").
		Where("o1.owner_id = ?", userID)
	errQuery := query.Find(&organizationBranchIDs).Error
	if errQuery != nil {
		if errors.Is(errQuery, gorm.ErrRecordNotFound) {
			return
		}
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errQuery)
	}
	return
}
