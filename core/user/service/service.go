// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/pkg/hash"
	userRepo "github.com/dedyf5/resik/repositories"
)

type Service struct {
	userRepo userRepo.IUser
	hasher   hash.IHash
	config   config.Config
}

func New(userRepo userRepo.IUser, hasher hash.IHash, config config.Config) *Service {
	return &Service{
		userRepo: userRepo,
		hasher:   hasher,
		config:   config,
	}
}
