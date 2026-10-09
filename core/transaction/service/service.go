// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/internal/identity"
	repo "github.com/dedyf5/resik/repositories"
)

type Service struct {
	config   config.Config
	resolver identity.IdentityResolver

	organizationRepo repo.IOrganization
	transactionRepo  repo.ITransaction
	userRepo         repo.IUser
}

func New(config config.Config, resolver identity.IdentityResolver, organizationRepo repo.IOrganization, transactionRepo repo.ITransaction, userRepo repo.IUser) *Service {
	return &Service{
		config:           config,
		resolver:         resolver,
		organizationRepo: organizationRepo,
		transactionRepo:  transactionRepo,
		userRepo:         userRepo,
	}
}
