// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	dtoMerchant "github.com/dedyf5/resik/core/merchant/dto"
	"github.com/dedyf5/resik/ctx"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	paramMerchant "github.com/dedyf5/resik/entities/merchant/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func (s *Service) MerchantInsert(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, status *resPkg.Status) {
	return s.merchantRepo.MerchantInsert(ctx, merchant)
}

func (s *Service) MerchantUpdate(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, status *resPkg.Status) {
	return s.merchantRepo.MerchantUpdate(ctx, merchant)
}

func (s *Service) MerchantListGet(param *paramMerchant.MerchantListGet) (res *dtoMerchant.MerchantList, status *resPkg.Status) {
	total, status := s.merchantRepo.MerchantListGetTotal(param)
	if status != nil {
		return nil, status
	}
	data, status := s.merchantRepo.MerchantListGetData(param)
	if status != nil {
		return nil, status
	}
	return &dtoMerchant.MerchantList{
		Data:  data,
		Total: total,
	}, nil
}

func (s *Service) MerchantDelete(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, status *resPkg.Status) {
	return s.merchantRepo.MerchantDelete(ctx, merchant)
}
