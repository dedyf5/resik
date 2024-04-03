// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	paramUser "github.com/dedyf5/resik/entities/user/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

//go:generate mockgen -source user.go -package mock -destination ./mock/user.go
type IService interface {
	Auth(param paramUser.Auth) (token string, status *resPkg.Status)
	AuthTokenGenerate(userID uint64, username string) (token string, status *resPkg.Status)
}
