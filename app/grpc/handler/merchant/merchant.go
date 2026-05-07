// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package merchant

import (
	merchantService "github.com/dedyf5/resik/core/merchant"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/internal/identity"
	validatorUtil "github.com/dedyf5/resik/utils/validator"
)

type MerchantHandler struct {
	log             *logCtx.Log
	validator       *validatorUtil.Validate
	resolver        identity.IdentityResolver
	merchantService merchantService.IService
}

func New(log *logCtx.Log, validator *validatorUtil.Validate, resolver identity.IdentityResolver, merchantService merchantService.IService) *MerchantHandler {
	return &MerchantHandler{
		log:             log,
		validator:       validator,
		resolver:        resolver,
		merchantService: merchantService,
	}
}

func (h *MerchantHandler) mustEmbedUnimplementedMerchantServiceServer() {}
