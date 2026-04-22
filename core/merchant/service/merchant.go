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

func (s *Service) MerchantGetByIDAndUserID(ctx *ctx.Ctx, merchantID, userID uint64) (merchant *merchantEntity.Merchant, err *resPkg.Status) {
	merchant, err = s.merchantRepo.MerchantGetByIDAndUserID(ctx, merchantID, userID)
	if merchant == nil || err != nil {
		return merchant, err
	}

	user, err := s.userRepo.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		localizer := ctx.Lang().Localizer
		return nil, resPkg.NewStatusMessage(
			http.StatusNotFound,
			term.NotFoundVal.Localize(localizer, term.User.Localize(localizer)),
			nil,
		)
	}

	merchant.User = *user

	return merchant, nil
}

func (s *Service) MerchantListGet(param *paramMerchant.MerchantListGet) (res *dtoMerchant.MerchantList, err *resPkg.Status) {
	total, err := s.merchantRepo.MerchantListGetTotal(param)
	if err != nil {
		return nil, err
	}
	data, err := s.merchantRepo.MerchantListGetData(param)
	if err != nil {
		return nil, err
	}
	return &dtoMerchant.MerchantList{
		Data:  data,
		Total: total,
	}, nil
}

func (s *Service) MerchantDelete(ctx *ctx.Ctx, merchant *merchantEntity.Merchant) (ok bool, err *resPkg.Status) {
	return s.merchantRepo.MerchantDelete(ctx, merchant)
}
