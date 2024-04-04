// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import "gorm.io/gorm"

type MerchantRepo struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *MerchantRepo {
	return &MerchantRepo{
		DB: DB,
	}
}
