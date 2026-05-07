// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"github.com/dedyf5/resik/config"
	trxService "github.com/dedyf5/resik/core/transaction"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/internal/identity"
	validatorUtil "github.com/dedyf5/resik/utils/validator"
)

type TransactionHandler struct {
	config     config.Config
	log        *logCtx.Log
	validator  *validatorUtil.Validate
	resolver   identity.IdentityResolver
	trxService trxService.IService
}

func New(config config.Config, log *logCtx.Log, validator *validatorUtil.Validate, resolver identity.IdentityResolver, trxService trxService.IService) *TransactionHandler {
	return &TransactionHandler{
		config:     config,
		log:        log,
		validator:  validator,
		resolver:   resolver,
		trxService: trxService,
	}
}

func (h *TransactionHandler) mustEmbedUnimplementedTransactionServiceServer() {}
