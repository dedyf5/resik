// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"gorm.io/gorm"
)

type TransactionRepo struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *TransactionRepo {
	return &TransactionRepo{
		DB: DB,
	}
}
