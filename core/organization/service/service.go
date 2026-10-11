// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/repositories"
)

type Service struct {
	config           config.Config
	organizationRepo repositories.IOrganization
	userRepo         repositories.IUser
}

func New(config config.Config, userRepo repositories.IUser, organizationRepo repositories.IOrganization) *Service {
	return &Service{
		config:           config,
		userRepo:         userRepo,
		organizationRepo: organizationRepo,
	}
}
