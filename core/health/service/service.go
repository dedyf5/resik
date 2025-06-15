// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"github.com/dedyf5/resik/core/health"
)

type Service struct {
	checkers []health.Checker
}

func New(checkers []health.Checker) health.IService {
	return &Service{
		checkers: checkers,
	}
}
