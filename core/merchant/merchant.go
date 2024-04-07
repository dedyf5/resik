// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	"github.com/dedyf5/resik/core/merchant/dto"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/merchant"
	"github.com/dedyf5/resik/entities/merchant/param"
	"github.com/dedyf5/resik/pkg/response"
)

//go:generate mockgen -source merchant.go -package mock -destination ./mock/merchant.go
type IService interface {
	MerchantInsert(ctx *ctx.Ctx, merchant *merchant.Merchant) (ok bool, status *response.Status)
	MerchantUpdate(ctx *ctx.Ctx, merchant *merchant.Merchant) (ok bool, status *response.Status)
	MerchantListGet(param *param.MerchantListGet) (res *dto.MerchantList, status *response.Status)
	MerchantDelete(ctx *ctx.Ctx, param *merchant.Merchant) (ok bool, status *response.Status)
}
