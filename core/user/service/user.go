// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"net/http"

	"github.com/dedyf5/resik/ctx"
	jwtCtx "github.com/dedyf5/resik/ctx/jwt"
	"github.com/dedyf5/resik/ctx/lang/term"
	paramUser "github.com/dedyf5/resik/entities/user/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

func (s *Service) Auth(param paramUser.Auth) (token string, err *resPkg.Status) {
	user, err := s.userRepo.UserByUsername(param.Ctx, param.Username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", resPkg.NewStatusMessage(
			http.StatusUnauthorized,
			term.IncorrectUsernameOrPassword.Localize(param.Ctx.Lang().Localizer),
			nil,
		)
	}

	if ok, err := s.hasher.Compare(param.Password, user.Password); !ok || err != nil {
		return "", resPkg.NewStatusMessage(
			http.StatusUnauthorized,
			term.IncorrectUsernameOrPassword.Localize(param.Ctx.Lang().Localizer),
			err,
		)
	}

	return s.AuthTokenGenerate(param.Ctx, user.ID, user.PublicID, user.Username)
}

func (s *Service) AuthTokenGenerate(ctx *ctx.Ctx, userID uint64, userPublicID uuidPkg.UUIDV7, username string) (token string, err *resPkg.Status) {
	merchantOutletIDs, err := s.userRepo.OutletMerchantByUserIDGetData(ctx, userID)
	if err != nil {
		return "", err
	}

	user := jwtCtx.User{
		Base: jwtCtx.Base{
			ID:       userID,
			PublicID: userPublicID,
		},
		Username: username,
	}

	merchantIDs, merchantPublicIDs, outletIDs, outletPublicIDs, err := merchantOutletIDs.UniqueIDs()
	if err != nil {
		return "", err
	}

	var merchants []jwtCtx.Base
	for i, v := range merchantIDs {
		merchants = append(merchants, jwtCtx.Base{
			ID:       v,
			PublicID: merchantPublicIDs[i],
		})
	}

	var outlets []jwtCtx.Base
	for i, v := range outletIDs {
		outlets = append(outlets, jwtCtx.Base{
			ID:       v,
			PublicID: outletPublicIDs[i],
		})
	}

	token, err = jwtCtx.AuthTokenGenerate(
		s.config.Module,
		s.config.Auth,
		user,
		merchants,
		outlets,
	)
	return
}
