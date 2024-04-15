// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"net/http"

	jwtCtx "github.com/dedyf5/resik/ctx/jwt"
	"github.com/dedyf5/resik/entities/outlet"
	paramUser "github.com/dedyf5/resik/entities/user/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
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

	return s.AuthTokenGenerate(user.ID, user.Username)
}

func (s *Service) AuthTokenGenerate(userID uint64, username string) (token string, status *resPkg.Status) {
	outlets, err := s.userRepo.OutletMerchantByUserIDGetData(userID)
	if err != nil {
		return "", err
	}

	merchantIDs, outletIDs := outlet.GetUniqueMerchantIDsAndOutletIDs(outlets)

	token, status = jwtCtx.AuthTokenGenerate(s.config.App, s.config.Auth, userID, username, merchantIDs, outletIDs)
	return
}
