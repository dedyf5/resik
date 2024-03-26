// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	userEntity "github.com/dedyf5/resik/entities/user"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func (s *Service) GetUserByUsernameAndPassword(userName string, password string) (user *userEntity.User, status *resPkg.Status) {
	return s.userRepo.GetUserByUsernameAndPassword(userName, password)
}
