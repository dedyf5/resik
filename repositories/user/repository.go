// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: DB,
	}
}
