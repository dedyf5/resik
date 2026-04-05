// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	repo "github.com/dedyf5/resik/repositories"
)

type Service struct {
	checkers []repo.ICheck
}

func New(checkers []repo.ICheck) *Service {
	return &Service{
		checkers: checkers,
	}
}
