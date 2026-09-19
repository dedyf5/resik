// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/internal/identity"
	"github.com/dedyf5/resik/repositories"
)

type Service struct {
	config           config.Config
	resolver         identity.IdentityResolver
	organizationRepo repositories.IOrganization
	userRepo         repositories.IUser
}

func New(config config.Config, resolver identity.IdentityResolver, userRepo repositories.IUser, organizationRepo repositories.IOrganization) *Service {
	return &Service{
		config:           config,
		resolver:         resolver,
		userRepo:         userRepo,
		organizationRepo: organizationRepo,
	}
}
