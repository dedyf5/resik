// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/config"
	userRepo "github.com/dedyf5/resik/repositories"
)

type Service struct {
	userRepo userRepo.IUser
	config   config.Config
}

func New(userRepo userRepo.IUser, config config.Config) *Service {
	return &Service{
		userRepo: userRepo,
		config:   config,
	}
}
