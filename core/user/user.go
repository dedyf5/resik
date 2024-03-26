// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	userEntity "github.com/dedyf5/resik/entities/user"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

type IService interface {
	GetUserByUsernameAndPassword(userName string, password string) (user *userEntity.User, status *resPkg.Status)
}
