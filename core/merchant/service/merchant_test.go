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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	configEntity "github.com/dedyf5/resik/entities/config"
	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	"github.com/dedyf5/resik/entities/merchant/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	merchantMock "github.com/dedyf5/resik/repositories/mock"
)

func TestMerchantInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	merchantRepo, ctx, merchantService := setup(ctrl)

	merchant := &merchantEntity.Merchant{}

	t.Run("MerchantInsert-ERROR", func(t *testing.T) {
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

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		okExpected := true
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantInsert(ctx, merchant).Return(okExpected, nil),
		)
		ok, err := merchantService.MerchantInsert(ctx, merchant)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})
}

func TestMerchantUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	merchantRepo, ctx, merchantService := setup(ctrl)

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

	merchantRepo, ctx, merchantService := setup(ctrl)

	merchant := &merchantEntity.Merchant{}

	t.Run("MerchantDelete-ERROR", func(t *testing.T) {
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

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		okExpected := true
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantDelete(ctx, merchant).Return(okExpected, nil),
		)
		ok, err := merchantService.MerchantDelete(ctx, merchant)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})
}

func TestMerchantListGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	merchantRepo, ctx, merchantService := setup(ctrl)

	param := &param.MerchantListGet{
		Ctx: ctx,
	}

	t.Run("MerchantListGetTotal-ERROR", func(t *testing.T) {
		var totalExpected uint64 = 0
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantListGetTotal(param).Return(totalExpected, statusErr),
		)
		res, err := merchantService.MerchantListGet(param)
		assert.Nil(t, res)
		assert.Equal(t, statusErr, err)
	})

	var totalExpected uint64 = 1
	t.Run("MerchantListGetData-ERROR", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantListGetTotal(param).Return(totalExpected, nil),
			merchantRepo.EXPECT().MerchantListGetData(param).Return(nil, statusErr),
		)
		res, err := merchantService.MerchantListGet(param)
		assert.Nil(t, res)
		assert.Equal(t, statusErr, err)
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		merchants := []merchantEntity.Merchant{
			{
				ID:   1,
				Name: "merchant1",
			},
		}
		gomock.InOrder(
			merchantRepo.EXPECT().MerchantListGetTotal(param).Return(totalExpected, nil),
			merchantRepo.EXPECT().MerchantListGetData(param).Return(merchants, nil),
		)
		res, err := merchantService.MerchantListGet(param)
		assert.Nil(t, err)
		assert.Equal(t, len(merchants), len(res.Data))
		assert.Equal(t, merchants[0].ID, res.Data[0].ID)
		assert.Equal(t, merchants[0].Name, res.Data[0].Name)
	})
}

func setup(ctrl *gomock.Controller) (merchantRepo *merchantMock.MockIMerchant, ctx *ctx.Ctx, merchantService *Service) {
	merchantRepo = merchantMock.NewMockIMerchant(ctrl)
	config, ctx := env()
	merchantService = New(merchantRepo, config)
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
	c = &ctx.Ctx{
		Context: context.Background(),
		Lang:    langCtx.NewLangTermDir(language.English, &language.English, "", fmt.Sprintf("%s%s", "../../../", langCtx.TermDir)),
	}
	return
}
