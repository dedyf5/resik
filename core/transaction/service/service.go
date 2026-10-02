// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/config"
	repo "github.com/dedyf5/resik/repositories"
)

type Service struct {
	config          config.Config
	transactionRepo repo.ITransaction
	userRepo        repo.IUser
}

func New(config config.Config, transactionRepo repo.ITransaction, userRepo repo.IUser) *Service {
	return &Service{
		config:          config,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}
