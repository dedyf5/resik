// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/config"
	orgRepo "github.com/dedyf5/resik/repositories/organization"
	userRepo "github.com/dedyf5/resik/repositories/user"
)

type Service struct {
	config           config.Config
	organizationRepo orgRepo.IOrganization
	userRepo         userRepo.IUser
}

func New(config config.Config, userRepo userRepo.IUser, organizationRepo orgRepo.IOrganization) *Service {
	return &Service{
		config:           config,
		userRepo:         userRepo,
		organizationRepo: organizationRepo,
	}
}
