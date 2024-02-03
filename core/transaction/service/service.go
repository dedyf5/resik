// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"time"

	"github.com/dedyf5/resik/config"
	statusCtx "github.com/dedyf5/resik/ctx/status"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	trxRepo "github.com/dedyf5/resik/repositories/transaction"
)

//go:generate mockgen -source service.go -package mock -destination ./mock/service.go
type IService interface {
	MerchantOmzetGet(param *paramTrx.MerchantOmzetGet) (res []trxEntity.MerchantOmzet, status *statusCtx.Status)
	OutletOmzet(outletID int64, date []time.Time) ([]trxEntity.OutletOmzet, error)
	GetUserByUserNameAndPassword(userName, password string) (*userEntity.User, error)
	ValidateAuthRequest(username, password string) error
	ValidateMerchantUser(merchantID, userID int64) (*merchantEntity.Merchant, error)
	ValidateOutletUser(outletID, createdBy int64) (*outletEntity.Outlet, error)
	Dates(date *time.Time, page int) []time.Time
}

type Service struct {
	transactionRepo trxRepo.ITransaction
	config          config.Config
}

func New(transactionRepo trxRepo.ITransaction, config config.Config) *Service {
	return &Service{
		transactionRepo: transactionRepo,
		config:          config,
	}
}
