// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/config"
	trxRepo "github.com/dedyf5/resik/repositories"
)

type Service struct {
	transactionRepo trxRepo.ITransaction
	config          config.Config
}

func New(transactionRepo trxRepo.ITransaction, config config.Config) *Service {
	return &Service{
		transactionRepo: transactionRepo,
		config:          config,
	}
}
