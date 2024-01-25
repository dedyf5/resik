// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"time"

	"github.com/dedyf5/resik/config"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	userEntity "github.com/dedyf5/resik/entities/user"
	trxRepo "github.com/dedyf5/resik/repositories/transaction"
)

//go:generate mockgen -source service.go -package mock -destination ./mock/service.go
type IService interface {
	MerchantOmzet(merchantID int64, date []time.Time) ([]trxEntity.MerchantOmzet, error)
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
