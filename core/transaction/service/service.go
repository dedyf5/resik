// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/internal/identity"
	orgRepo "github.com/dedyf5/resik/repositories/organization"
	trxRepo "github.com/dedyf5/resik/repositories/transaction"
	userRepo "github.com/dedyf5/resik/repositories/user"
)

type Service struct {
	config   config.Config
	resolver identity.IdentityResolver

	organizationRepo orgRepo.IOrganization
	transactionRepo  trxRepo.ITransaction
	userRepo         userRepo.IUser
}

func New(config config.Config, resolver identity.IdentityResolver, organizationRepo orgRepo.IOrganization, transactionRepo trxRepo.ITransaction, userRepo userRepo.IUser) *Service {
	return &Service{
		config:           config,
		resolver:         resolver,
		organizationRepo: organizationRepo,
		transactionRepo:  transactionRepo,
		userRepo:         userRepo,
	}
}
