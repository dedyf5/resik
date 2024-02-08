package transaction

import (
	"time"

	trxDTO "github.com/dedyf5/resik/core/transaction/dto"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

//go:generate mockgen -source transaction.go -package mock -destination ./mock/transaction.go
type IService interface {
	MerchantOmzetGet(param *paramTrx.MerchantOmzetGet) (res *trxDTO.MerchantOmzet, status *resPkg.Status)
	OutletOmzetGet(param *paramTrx.OutletOmzetGet) (res *trxDTO.OutletOmzet, status *resPkg.Status)
	GetUserByUserNameAndPassword(userName, password string) (*userEntity.User, error)
	ValidateAuthRequest(username, password string) error
	ValidateMerchantUser(merchantID, userID int64) (*merchantEntity.Merchant, error)
	ValidateOutletUser(outletID, createdBy int64) (*outletEntity.Outlet, error)
	Dates(date *time.Time, page int) []time.Time
}
