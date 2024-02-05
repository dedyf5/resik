package transaction

import (
	"time"

	trxDTO "github.com/dedyf5/resik/core/transaction/dto"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	statusPkg "github.com/dedyf5/resik/pkg/status"
)

//go:generate mockgen -source transaction.go -package mock -destination ./mock/transaction.go
type IService interface {
	MerchantOmzetGet(param *paramTrx.MerchantOmzetGet) (res *trxDTO.MerchantOmzet, status *statusPkg.Status)
	OutletOmzet(outletID int64, date []time.Time) ([]trxEntity.OutletOmzet, error)
	GetUserByUserNameAndPassword(userName, password string) (*userEntity.User, error)
	ValidateAuthRequest(username, password string) error
	ValidateMerchantUser(merchantID, userID int64) (*merchantEntity.Merchant, error)
	ValidateOutletUser(outletID, createdBy int64) (*outletEntity.Outlet, error)
	Dates(date *time.Time, page int) []time.Time
}
