// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"github.com/dedyf5/resik/config"
	trxService "github.com/dedyf5/resik/core/transaction"
	logCtx "github.com/dedyf5/resik/ctx/log"
	validatorUtil "github.com/dedyf5/resik/utils/validator"
)

type TransactionHandler struct {
	config     config.Config
	log        *logCtx.Log
	validator  *validatorUtil.Validate
	trxService trxService.IService
}

func New(config config.Config, log *logCtx.Log, validator *validatorUtil.Validate, trxService trxService.IService) *TransactionHandler {
	return &TransactionHandler{
		config:     config,
		log:        log,
		validator:  validator,
		trxService: trxService,
	}
}

func (h *TransactionHandler) mustEmbedUnimplementedTransactionServiceServer() {}
