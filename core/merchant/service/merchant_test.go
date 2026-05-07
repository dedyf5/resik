// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/dedyf5/resik/config"
	dtoMerchant "github.com/dedyf5/resik/core/merchant/dto"
	"github.com/dedyf5/resik/ctx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/log"
	configEntity "github.com/dedyf5/resik/entities/config"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	"github.com/dedyf5/resik/entities/merchant/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	identityMock "github.com/dedyf5/resik/internal/identity/mock"
	resPkg "github.com/dedyf5/resik/pkg/response"
	repoMock "github.com/dedyf5/resik/repositories/mock"
)

func TestMerchantInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resolver, ctx, merchantRepo, _, merchantService := setup(ctrl)

	userID := ctx.UserClaims().UserID()
	merchant := &merchantEntity.Merchant{}

	t.Run("MerchantInsert-ERROR1", func(t *testing.T) {
		okExpected := false
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantInsert(ctx, merchant).Return(okExpected, statusErr),
		)
		ok, err := merchantService.MerchantInsert(ctx, merchant)
		assert.Equal(t, okExpected, ok)
		assert.Equal(t, statusErr, err)
	})

	t.Run("MerchantInsert-ERROR2", func(t *testing.T) {
		okExpected := true
		errExpected := errors.New("ERROR")
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantInsert(ctx, merchant).Return(okExpected, nil),
			resolver.EXPECT().InvalidateUserAccessMerchant(ctx.Context, userID).Return(errExpected),
		)
		ok, err := merchantService.MerchantInsert(ctx, merchant)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		okExpected := true
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantInsert(ctx, merchant).Return(okExpected, nil),
			resolver.EXPECT().InvalidateUserAccessMerchant(ctx.Context, userID).Return(nil),
		)
		ok, err := merchantService.MerchantInsert(ctx, merchant)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})
}

func TestMerchantUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, ctx, merchantRepo, _, merchantService := setup(ctrl)

	merchant := &merchantEntity.Merchant{}

	t.Run("MerchantUpdate-ERROR", func(t *testing.T) {
		okExpected := false
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantUpdate(ctx, merchant).Return(okExpected, statusErr),
		)
		ok, err := merchantService.MerchantUpdate(ctx, merchant)
		assert.Equal(t, okExpected, ok)
		assert.Equal(t, statusErr, err)
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		okExpected := true
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantUpdate(ctx, merchant).Return(true, nil),
		)
		ok, err := merchantService.MerchantUpdate(ctx, merchant)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})
}

func TestMerchantDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resolver, ctx, merchantRepo, _, merchantService := setup(ctrl)

	userID := ctx.UserClaims().UserID()
	merchant := &merchantEntity.Merchant{}

	t.Run("MerchantDelete-ERROR1", func(t *testing.T) {
		okExpected := false
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantDelete(ctx, merchant).Return(okExpected, statusErr),
		)
		ok, err := merchantService.MerchantDelete(ctx, merchant)
		assert.Equal(t, okExpected, ok)
		assert.Equal(t, statusErr, err)
	})

	t.Run("MerchantDelete-ERROR2", func(t *testing.T) {
		okExpected := true
		errExpected := errors.New("ERROR")
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantDelete(ctx, merchant).Return(okExpected, nil),
			resolver.EXPECT().InvalidateUserAccessMerchant(ctx.Context, userID).Return(errExpected),
		)
		ok, err := merchantService.MerchantDelete(ctx, merchant)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		okExpected := true
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantDelete(ctx, merchant).Return(okExpected, nil),
			resolver.EXPECT().InvalidateUserAccessMerchant(ctx.Context, userID).Return(nil),
		)
		ok, err := merchantService.MerchantDelete(ctx, merchant)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})
}

func TestMerchantGetByIDAndOwnerID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, ctx, merchantRepo, userRepo, merchantService := setup(ctrl)

	merchant := &merchantEntity.Merchant{}
	user := &userEntity.User{}
	users := userEntity.Users{
		*user,
	}
	userIDs := []uint64{user.ID}

	t.Run("MerchantGetByIDAndOwnerID-ERROR MerchantGetByIDAndOwnerID", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantGetByIDAndOwnerID(ctx, merchant.ID, user.ID).Return(merchant, statusErr),
		)
		_, err := merchantService.MerchantGetByIDAndOwnerID(ctx, merchant.ID, user.ID)
		assert.Equal(t, statusErr, err)
	})

	t.Run("MerchantGetByIDAndOwnerID-ERROR UsersGetByIDs ERROR", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantGetByIDAndOwnerID(ctx, merchant.ID, user.ID).Return(merchant, nil),
			userRepo.EXPECT().UsersGetByIDs(ctx, userIDs).Return(users, statusErr),
		)
		_, err := merchantService.MerchantGetByIDAndOwnerID(ctx, merchant.ID, user.ID)
		assert.Equal(t, statusErr, err)
	})

	t.Run("MerchantGetByIDAndOwnerID-ERROR UsersGetByIDs Not Found", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusNotFound,
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantGetByIDAndOwnerID(ctx, merchant.ID, user.ID).Return(merchant, nil),
			userRepo.EXPECT().UsersGetByIDs(ctx, userIDs).Return(nil, nil),
		)
		_, err := merchantService.MerchantGetByIDAndOwnerID(ctx, merchant.ID, user.ID)
		assert.Equal(t, statusErr.Code, err.Code)
	})

	t.Run("ALL-EMPTY", func(t *testing.T) {
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantGetByIDAndOwnerID(ctx, merchant.ID, user.ID).Return(nil, nil),
		)
		res, err := merchantService.MerchantGetByIDAndOwnerID(ctx, merchant.ID, user.ID)
		assert.Nil(t, err)
		assert.Nil(t, res)
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantGetByIDAndOwnerID(ctx, merchant.ID, user.ID).Return(merchant, nil),
			userRepo.EXPECT().UsersGetByIDs(ctx, userIDs).Return(users, nil),
		)
		res, err := merchantService.MerchantGetByIDAndOwnerID(ctx, merchant.ID, user.ID)
		assert.Nil(t, err)
		assert.Equal(t, merchant.PublicID, res.PublicID)
	})
}

func TestMerchantsGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, ctx, merchantRepo, userRepo, merchantService := setup(ctrl)

	param := &param.MerchantsGet{
		Ctx: ctx,
	}

	merchants := merchantEntity.Merchants{
		{
			ID:   1,
			Name: "merchant1",
		},
	}

	user := &userEntity.User{}
	users := userEntity.Users{
		*user,
	}
	userIDs := []uint64{user.ID}

	t.Run("MerchantListGetTotal-ERROR", func(t *testing.T) {
		var totalExpected int64 = 0
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantsGetTotal(param).Return(totalExpected, statusErr),
		)
		res, err := merchantService.MerchantsGet(param)
		assert.Nil(t, res)
		assert.Equal(t, statusErr, err)
	})

	t.Run("MerchantListGetTotal-0", func(t *testing.T) {
		var totalExpected int64 = 0
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantsGetTotal(param).Return(totalExpected, nil),
		)
		res, err := merchantService.MerchantsGet(param)
		assert.Nil(t, err)
		assert.Equal(t, dtoMerchant.MerchatsResultEmpty, *res)
	})

	var totalExpected int64 = 1
	t.Run("MerchantListGetData-ERROR", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantsGetTotal(param).Return(totalExpected, nil),
			merchantRepo.EXPECT().MerchantsGetData(param).Return(nil, statusErr),
		)
		res, err := merchantService.MerchantsGet(param)
		assert.Nil(t, res)
		assert.Equal(t, statusErr, err)
	})

	t.Run("MerchantListGetData-0", func(t *testing.T) {
		merchantsEmpty := merchantEntity.Merchants{}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantsGetTotal(param).Return(totalExpected, nil),
			merchantRepo.EXPECT().MerchantsGetData(param).Return(merchantsEmpty, nil),
		)
		res, err := merchantService.MerchantsGet(param)
		assert.Nil(t, err)
		assert.Equal(t, dtoMerchant.MerchatsResultEmpty, *res)
	})

	t.Run("UsersGetByIDs-ERROR", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantsGetTotal(param).Return(totalExpected, nil),
			merchantRepo.EXPECT().MerchantsGetData(param).Return(merchants, nil),
			userRepo.EXPECT().UsersGetByIDs(ctx, userIDs).Return(nil, statusErr),
		)
		res, err := merchantService.MerchantsGet(param)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.Code, err.Code)
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantsGetTotal(param).Return(totalExpected, nil),
			merchantRepo.EXPECT().MerchantsGetData(param).Return(merchants, nil),
			userRepo.EXPECT().UsersGetByIDs(ctx, userIDs).Return(users, nil),
		)
		res, err := merchantService.MerchantsGet(param)
		assert.Nil(t, err)
		assert.Len(t, res.Data, len(merchants))
		assert.Equal(t, merchants[0].ID, res.Data[0].ID)
		assert.Equal(t, merchants[0].Name, res.Data[0].Name)
	})
}

func setup(ctrl *gomock.Controller) (resolver *identityMock.MockIdentityResolver, ctx *ctx.Ctx, merchantRepo *repoMock.MockIMerchant, userRepo *repoMock.MockIUser, merchantService *Service) {
	resolver = identityMock.NewMockIdentityResolver(ctrl)
	merchantRepo = repoMock.NewMockIMerchant(ctrl)
	userRepo = repoMock.NewMockIUser(ctrl)
	config, ctx := env()
	merchantService = New(config, resolver, userRepo, merchantRepo)
	return
}

func env() (conf config.Config, c *ctx.Ctx) {
	conf = config.Config{
		Module: configEntity.Module{
			Name:        "REST",
			NameKey:     "rest",
			LangDefault: language.English,
			Host:        "0.0.0.0",
			Port:        8081,
		},
	}
	context := context.WithValue(context.Background(), langCtx.ContextKey, langCtx.NewLangLocaleDir(language.English, &language.English, "", fmt.Sprintf("%s%s", "../../../", langCtx.LocaleDir)))
	c, _ = ctx.NewCtx(context, log.Get(configEntity.Log{}, conf.Module))
	return
}
