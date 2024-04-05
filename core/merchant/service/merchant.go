// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/ctx"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func (s *Service) MerchantInsert(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, status *resPkg.Status) {
	return s.merchantRepo.MerchantInsert(ctx, merchant)
}

func (s *Service) MerchantUpdate(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, status *resPkg.Status) {
	return s.merchantRepo.MerchantUpdate(ctx, merchant)
}
