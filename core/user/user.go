// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	paramUser "github.com/dedyf5/resik/entities/user/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

type IService interface {
	Auth(param paramUser.Auth) (token string, status *resPkg.Status)
}
