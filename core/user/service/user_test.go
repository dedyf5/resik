// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/ctx"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/log"
	configEntity "github.com/dedyf5/resik/entities/config"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	userEntity "github.com/dedyf5/resik/entities/user"
	paramUser "github.com/dedyf5/resik/entities/user/param"
	hashMock "github.com/dedyf5/resik/pkg/hash/mock"
	resPkg "github.com/dedyf5/resik/pkg/response"
	userMock "github.com/dedyf5/resik/repositories/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"
)

const (
	userID   uint64 = 1
	username        = "admin1"
)

func TestAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo, hasher, ctx, userService := setup(ctrl)

	param := paramUser.Auth{
		Ctx:      ctx,
		Username: username,
		Password: "",
	}

	userExpected := &userEntity.User{
		ID:       1,
		Username: username,
	}

	t.Run("UserByUsername-ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			userRepo.EXPECT().UserByUsername(param.Ctx, param.Username).Return(nil, statusErr),
		)
		token, err := userService.Auth(param)
		assert.NotNil(t, err)
		assert.Empty(t, token)
	})

	t.Run("UserByUsername-ERROR-401-1", func(t *testing.T) {
		gomock.InOrder(
			userRepo.EXPECT().UserByUsername(param.Ctx, param.Username).Return(nil, nil),
		)
		statusErr := &resPkg.Status{
			Code:    http.StatusUnauthorized,
			Message: param.Ctx.Lang().GetByMessageID("incorrect_username_or_password"),
		}
		token, err := userService.Auth(param)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr, err)
		assert.Empty(t, token)
	})

	t.Run("UserByUsername-ERROR-401-2", func(t *testing.T) {
		gomock.InOrder(
			userRepo.EXPECT().UserByUsername(param.Ctx, param.Username).Return(userExpected, nil),
			hasher.EXPECT().Compare(param.Password, userExpected.Password).Return(false, nil),
		)
		statusErr := &resPkg.Status{
			Code:    http.StatusUnauthorized,
			Message: param.Ctx.Lang().GetByMessageID("incorrect_username_or_password"),
		}
		token, err := userService.Auth(param)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr, err)
		assert.Empty(t, token)
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		userExpected := &userEntity.User{
			ID:       1,
			Username: username,
		}
		gomock.InOrder(
			userRepo.EXPECT().UserByUsername(param.Ctx, param.Username).Return(userExpected, nil),
			hasher.EXPECT().Compare(param.Password, userExpected.Password).Return(true, nil),
			userRepo.EXPECT().OutletMerchantByUserIDGetData(userID).Return(outletsExpected(), nil),
		)
		token, err := userService.Auth(param)
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestAuthTokenGenerate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo, _, _, userService := setup(ctrl)

	t.Run("OutletMerchantByUserIDGetData-ERROR", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			userRepo.EXPECT().OutletMerchantByUserIDGetData(userID).Return(nil, statusErr),
		)
		res, err := userService.AuthTokenGenerate(userID, username)
		assert.Empty(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr, err)
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		gomock.InOrder(
			userRepo.EXPECT().OutletMerchantByUserIDGetData(userID).Return(outletsExpected(), nil),
		)
		res, err := userService.AuthTokenGenerate(userID, username)
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
	})
}

func outletsExpected() []outletEntity.Outlet {
	return []outletEntity.Outlet{
		{
			ID:         1,
			MerchantID: 1,
		},
		{
			ID:         3,
			MerchantID: 1,
		},
		{
			MerchantID: 3,
		},
	}
}

func setup(ctrl *gomock.Controller) (userRepo *userMock.MockIUser, hasher *hashMock.MockIHash, ctx *ctx.Ctx, userService *Service) {
	userRepo = userMock.NewMockIUser(ctrl)
	hasher = hashMock.NewMockIHash(ctrl)
	config, ctx := env()
	userService = New(userRepo, hasher, config)
	return
}

func env() (conf config.Config, c *ctx.Ctx) {
	conf = config.Config{
		App: configEntity.App{
			Name:        "Resik",
			Version:     "0.1",
			LangDefault: language.English,
			Host:        "0.0.0.0",
			Port:        8081,
		},
	}
	context := context.WithValue(context.Background(), langCtx.ContextKey, langCtx.NewLangTermDir(language.English, &language.English, "", fmt.Sprintf("%s%s", "../../../", langCtx.TermDir)))
	c, _ = ctx.NewCtx(context, &log.Log{})
	return
}
