// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	merchantService "github.com/dedyf5/resik/core/merchant"
	logCtx "github.com/dedyf5/resik/ctx/log"
	validatorUtil "github.com/dedyf5/resik/utils/validator"
)

type MerchantHandler struct {
	log             *logCtx.Log
	validator       *validatorUtil.Validate
	merchantService merchantService.IService
}

func New(log *logCtx.Log, validator *validatorUtil.Validate, merchantService merchantService.IService) *MerchantHandler {
	return &MerchantHandler{
		log:             log,
		validator:       validator,
		merchantService: merchantService,
	}
}

func (h *MerchantHandler) mustEmbedUnimplementedMerchantServiceServer() {}
