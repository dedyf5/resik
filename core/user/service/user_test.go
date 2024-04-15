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
	configEntity "github.com/dedyf5/resik/entities/config"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	userEntity "github.com/dedyf5/resik/entities/user"
	paramUser "github.com/dedyf5/resik/entities/user/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	userMock "github.com/dedyf5/resik/repositories/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

const (
	userID   uint64 = 1
	username        = "admin1"
)

func TestAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo, ctx, userService := setup(ctrl)

	param := paramUser.Auth{
		Ctx:      ctx,
		Username: username,
		Password: "",
	}

	t.Run("UserByUsernameAndPasswordGetData-ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			userRepo.EXPECT().UserByUsernameAndPasswordGetData(param).Return(nil, statusErr),
		)
		token, err := userService.Auth(param)
		assert.NotNil(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("UserByUsernameAndPasswordGetData-ERROR-401", func(t *testing.T) {
		gomock.InOrder(
			userRepo.EXPECT().UserByUsernameAndPasswordGetData(param).Return(nil, nil),
		)
		statusErr := &resPkg.Status{
			Code:    http.StatusUnauthorized,
			Message: param.Ctx.Lang.GetByMessageID("incorrect_username_or_password"),
		}
		token, err := userService.Auth(param)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr, err)
		assert.Equal(t, "", token)
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		userExpected := &userEntity.User{
			ID:       1,
			Username: username,
		}
		gomock.InOrder(
			userRepo.EXPECT().UserByUsernameAndPasswordGetData(param).Return(userExpected, nil),
			userRepo.EXPECT().OutletMerchantByUserIDGetData(userID).Return(outletsExpected(), nil),
		)
		token, err := userService.Auth(param)
		assert.Nil(t, err)
		assert.NotEqual(t, "", token)
	})
}

func TestAuthTokenGenerate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo, _, userService := setup(ctrl)

	t.Run("OutletMerchantByUserIDGetData-ERROR", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			userRepo.EXPECT().OutletMerchantByUserIDGetData(userID).Return(nil, statusErr),
		)
		res, err := userService.AuthTokenGenerate(userID, username)
		assert.Equal(t, "", res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr, err)
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		gomock.InOrder(
			userRepo.EXPECT().OutletMerchantByUserIDGetData(userID).Return(outletsExpected(), nil),
		)
		res, err := userService.AuthTokenGenerate(userID, username)
		assert.Nil(t, err)
		assert.NotEqual(t, "", res)
	})
}

func setup(ctrl *gomock.Controller) (userRepo *userMock.MockIUser, ctx *ctx.Ctx, userService *Service) {
	userRepo = userMock.NewMockIUser(ctrl)
	config, ctx := env()
	userService = New(userRepo, config)
	return
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
	c = &ctx.Ctx{
		Context: context.Background(),
		Lang:    langCtx.NewLangTermDir(language.English, &language.English, "", fmt.Sprintf("%s%s", "../../../", langCtx.TermDir)),
	}
	return
}
