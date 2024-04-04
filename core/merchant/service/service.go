// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/repositories"
)

type Service struct {
	merchantRepo repositories.IMerchant
	config       config.Config
}

func New(merchantRepo repositories.IMerchant, config config.Config) *Service {
	return &Service{
		merchantRepo: merchantRepo,
		config:       config,
	}
}
