// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package mysql

import "gorm.io/gorm"

type OrganizationRepo struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *OrganizationRepo {
	return &OrganizationRepo{
		DB: db,
	}
}
