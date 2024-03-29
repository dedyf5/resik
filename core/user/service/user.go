// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"net/http"

	paramUser "github.com/dedyf5/resik/entities/user/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	jwtUtil "github.com/dedyf5/resik/utils/jwt"
)

func (s *Service) Auth(param paramUser.Auth) (token string, status *resPkg.Status) {
	user, err := s.userRepo.UserByUsernameAndPasswordGetData(param)
	if err != nil {
		return "", err
	}
	if user == nil && err == nil {
		return "", &resPkg.Status{
			Code:    http.StatusBadRequest,
			Message: param.Ctx.Lang.GetByMessageID("incorrect_username_or_password"),
		}
	}

	merchantIDs, err := s.userRepo.MerchantIDsByUserIDGetData(user.ID)
	if err != nil {
		return "", &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: err,
		}
	}

	token, status = jwtUtil.AuthTokenGenerate(s.config.App, s.config.Auth, user.ID, user.Username.String, merchantIDs)
	return
}
