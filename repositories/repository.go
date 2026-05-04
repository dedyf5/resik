// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package repositories

import (
	"github.com/dedyf5/resik/ctx"
	checkEntity "github.com/dedyf5/resik/entities/check"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	paramMerchant "github.com/dedyf5/resik/entities/merchant/param"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

//go:generate mockgen -source repository.go -package mock -destination ./mock/repository.go
type ICheck interface {
	Check() checkEntity.CheckDetail
}

type ITransaction interface {
	MerchantOmzetGetData(param *paramTrx.MerchantOmzetGet) (res []trxEntity.MerchantOmzet, err *resPkg.Status)
	MerchantOmzetGetTotal(param *paramTrx.MerchantOmzetGet) (total int64, err *resPkg.Status)
	OutletOmzetGetData(param *paramTrx.OutletOmzetGet) (res []trxEntity.OutletOmzet, err *resPkg.Status)
	OutletOmzetGetTotal(param *paramTrx.OutletOmzetGet) (total int64, err *resPkg.Status)
}

type IUser interface {
	UserByID(ctx *ctx.Ctx, userID uint64) (user *userEntity.User, err *resPkg.Status)
	UserByUsername(ctx *ctx.Ctx, username string) (user *userEntity.User, err *resPkg.Status)
	UsersGetByIDs(ctx *ctx.Ctx, userIDs []uint64) (users userEntity.Users, err *resPkg.Status)
	MerchantIDsByUserIDGetData(userID uint64) (merchantIDs []uint64, err *resPkg.Status)
	OutletMerchantByUserIDGetData(ctx *ctx.Ctx, userID uint64) (merchantOutletIDs userEntity.MerchantOutletIDs, err *resPkg.Status)
}

type IMerchant interface {
	MerchantInsert(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, err *resPkg.Status)
	MerchantUpdate(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, err *resPkg.Status)
	MerchantGetByIDAndOwnerID(ctx *ctx.Ctx, merchantID, ownerID uint64) (merchant *merchantEntity.Merchant, err *resPkg.Status)
	MerchantsGetData(param *paramMerchant.MerchantsGet) (merchants merchantEntity.Merchants, err *resPkg.Status)
	MerchantsGetTotal(param *paramMerchant.MerchantsGet) (total int64, err *resPkg.Status)
	MerchantDelete(c *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, err *resPkg.Status)
}
