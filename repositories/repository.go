// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package repositories

import (
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

//go:generate mockgen -source repository.go -package mock -destination ./mock/repository.go
type ITransaction interface {
	MerchantOmzetGetData(param *paramTrx.MerchantOmzetGet) (res []trxEntity.MerchantOmzet, status *resPkg.Status)
	MerchantOmzetGetTotal(param *paramTrx.MerchantOmzetGet) (total uint64, status *resPkg.Status)
	OutletOmzetGetData(param *paramTrx.OutletOmzetGet) (res []trxEntity.OutletOmzet, status *resPkg.Status)
	OutletOmzetGetTotal(param *paramTrx.OutletOmzetGet) (total uint64, status *resPkg.Status)
	GetMerchantByID(merchantID int64) (*merchantEntity.Merchant, error)
	GetMerchantByIDAndUserID(merchantID int64, userID int64) (*merchantEntity.Merchant, error)
	GetOutletByID(outletID int64) (*outletEntity.Outlet, error)
	GetOutletByIDAndCreatedBy(outletID int64, createdBy int64) (*outletEntity.Outlet, error)
}

type IUser interface {
	GetUserByUsernameAndPassword(userName string, password string) (user *userEntity.User, status *resPkg.Status)
}
