// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package repositories

import (
	"github.com/dedyf5/resik/ctx"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	paramMerchant "github.com/dedyf5/resik/entities/merchant/param"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	paramUser "github.com/dedyf5/resik/entities/user/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

//go:generate mockgen -source repository.go -package mock -destination ./mock/repository.go
type ITransaction interface {
	MerchantOmzetGetData(param *paramTrx.MerchantOmzetGet) (res []trxEntity.MerchantOmzet, status *resPkg.Status)
	MerchantOmzetGetTotal(param *paramTrx.MerchantOmzetGet) (total uint64, status *resPkg.Status)
	OutletOmzetGetData(param *paramTrx.OutletOmzetGet) (res []trxEntity.OutletOmzet, status *resPkg.Status)
	OutletOmzetGetTotal(param *paramTrx.OutletOmzetGet) (total uint64, status *resPkg.Status)
	GetMerchantByID(merchantID uint64) (*merchantEntity.Merchant, error)
	GetMerchantByIDAndUserID(merchantID uint64, userID uint64) (*merchantEntity.Merchant, error)
	GetOutletByID(outletID uint64) (*outletEntity.Outlet, error)
	GetOutletByIDAndCreatedBy(outletID uint64, createdBy uint64) (*outletEntity.Outlet, error)
}

type IUser interface {
	UserByUsernameAndPasswordGetData(param paramUser.Auth) (user *userEntity.User, status *resPkg.Status)
	MerchantIDsByUserIDGetData(userID uint64) (merchantIDs []uint64, status *resPkg.Status)
	OutletMerchantByUserIDGetData(userID uint64) (outlets []outletEntity.Outlet, status *resPkg.Status)
}

type IMerchant interface {
	MerchantInsert(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, status *resPkg.Status)
	MerchantUpdate(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, status *resPkg.Status)
	MerchantListGetData(param *paramMerchant.MerchantListGet) (merchant []merchantEntity.Merchant, status *resPkg.Status)
	MerchantListGetTotal(param *paramMerchant.MerchantListGet) (total uint64, status *resPkg.Status)
}
