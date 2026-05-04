// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"net/http"

	dtoMerchant "github.com/dedyf5/resik/core/merchant/dto"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	paramMerchant "github.com/dedyf5/resik/entities/merchant/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func (s *Service) MerchantInsert(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, err *resPkg.Status) {
	return s.merchantRepo.MerchantInsert(ctx, merchant)
}

func (s *Service) MerchantUpdate(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, err *resPkg.Status) {
	return s.merchantRepo.MerchantUpdate(ctx, merchant)
}

func (s *Service) MerchantGetByIDAndOwnerID(ctx *ctx.Ctx, merchantID, ownerID uint64) (*dtoMerchant.Merchant, *resPkg.Status) {
	merchant, err := s.merchantRepo.MerchantGetByIDAndOwnerID(ctx, merchantID, ownerID)
	if err != nil {
		return nil, err
	}

	if merchant == nil {
		return nil, nil
	}

	users, err := s.userRepo.UsersGetByIDs(ctx, merchant.UniqueAllUserIDs())
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		localizer := ctx.Lang().Localizer
		return nil, resPkg.NewStatusMessage(
			http.StatusNotFound,
			term.NotFoundVal.Localize(localizer, term.User.Localize(localizer)),
			nil,
		)
	}

	res := dtoMerchant.MerchantFromEntity(*merchant, users.UniqueMap())
	return &res, nil
}

func (s *Service) MerchantsGet(param *paramMerchant.MerchantsGet) (res *dtoMerchant.MerchantsResult, err *resPkg.Status) {
	total, err := s.merchantRepo.MerchantsGetTotal(param)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &dtoMerchant.MerchatsResultEmpty, nil
	}

	merchants, err := s.merchantRepo.MerchantsGetData(param)
	if err != nil {
		return nil, err
	}

	if len(merchants) == 0 {
		return &dtoMerchant.MerchatsResultEmpty, nil
	}

	users, err := s.userRepo.UsersGetByIDs(param.Ctx, merchants.UniqueAllUserIDs())
	if err != nil {
		return nil, err
	}

	return &dtoMerchant.MerchantsResult{
		Data:  dtoMerchant.MerchantsFromEntity(merchants, users.UniqueMap()),
		Total: total,
	}, nil
}

func (s *Service) MerchantDelete(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, err *resPkg.Status) {
	return s.merchantRepo.MerchantDelete(ctx, merchant)
}
