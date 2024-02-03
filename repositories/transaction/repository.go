// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"time"

	statusCtx "github.com/dedyf5/resik/ctx/status"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	"gorm.io/gorm"
)

//go:generate mockgen -source repository.go -package mock -destination ./mock/repository.go
type ITransaction interface {
	MerchantOmzetGet(param *paramTrx.MerchantOmzetGet) (res []trxEntity.MerchantOmzet, status *statusCtx.Status)
	OutletOmzet(outletID int64, date time.Time) (*trxEntity.OutletOmzet, error)
	GetMerchantByID(merchantID int64) (*merchantEntity.Merchant, error)
	GetMerchantByIDAndUserID(merchantID int64, userID int64) (*merchantEntity.Merchant, error)
	GetOutletByID(outletID int64) (*outletEntity.Outlet, error)
	GetOutletByIDAndCreatedBy(outletID int64, createdBy int64) (*outletEntity.Outlet, error)
	GetUserByUserNameAndPassword(userName string, password string) (*userEntity.User, error)
}

type TransactionRepo struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) *TransactionRepo {
	return &TransactionRepo{
		DB: DB,
	}
}
