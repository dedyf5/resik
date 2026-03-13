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
	"github.com/dedyf5/resik/ctx"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/log"
	configEntity "github.com/dedyf5/resik/entities/config"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	trxParam "github.com/dedyf5/resik/entities/transaction/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	trxRepoMock "github.com/dedyf5/resik/repositories/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"
)

func TestMerchantOmzetGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trxRepo, ctx, trxService := setup(ctrl)

	p := trxParam.MerchantOmzetGet{
		Ctx:        ctx,
		MerchantID: 1,
	}

	t.Run("MerchantOmzetGetTotal-ERROR", func(t *testing.T) {
		errNative := errors.New("failed to get total")
		statusErr := &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errNative,
		}
		gomock.InOrder(
			trxRepo.EXPECT().MerchantOmzetGetTotal(&p).Return(int64(0), statusErr),
		)
		res, err := trxService.MerchantOmzetGet(&p)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.MessageOrDefault(), err.MessageOrDefault())
		assert.Equal(t, statusErr.CauseError.Error(), err.CauseError.Error())
	})

	t.Run("MerchantOmzetGetData-ERROR", func(t *testing.T) {
		errNative := errors.New("failed to get data")
		statusErr := &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errNative,
		}
		gomock.InOrder(
			trxRepo.EXPECT().MerchantOmzetGetTotal(&p).Return(int64(1), nil),
			trxRepo.EXPECT().MerchantOmzetGetData(&p).Return(nil, statusErr),
		)
		res, err := trxService.MerchantOmzetGet(&p)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.MessageOrDefault(), err.MessageOrDefault())
		assert.Equal(t, statusErr.CauseError.Error(), err.CauseError.Error())
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		expRes := make([]trxEntity.MerchantOmzet, 0, 1)
		expRes = append(expRes, trxEntity.MerchantOmzet{
			MerchantID:   1,
			MerchantName: "Merchant Name",
			Omzet:        500.75,
			Period:       "2024-03-07",
		})
		resInt64 := int64(len(expRes))
		gomock.InOrder(
			trxRepo.EXPECT().MerchantOmzetGetTotal(&p).Return(resInt64, nil),
			trxRepo.EXPECT().MerchantOmzetGetData(&p).Return(expRes, nil),
		)
		res, err := trxService.MerchantOmzetGet(&p)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, resInt64, res.Total)
		assert.Equal(t, int(resInt64), len(res.Data))
		assert.Equal(t, expRes, res.Data)
	})
}

func TestOutletOmzetGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trxRepo, ctx, trxService := setup(ctrl)

	p := trxParam.OutletOmzetGet{
		Ctx:      ctx,
		OutletID: 1,
	}

	t.Run("OutletOmzetGetTotal-ERROR", func(t *testing.T) {
		errNative := errors.New("failed to get total")
		statusErr := &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errNative,
		}
		gomock.InOrder(
			trxRepo.EXPECT().OutletOmzetGetTotal(&p).Return(int64(0), statusErr),
		)
		res, err := trxService.OutletOmzetGet(&p)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.MessageOrDefault(), err.MessageOrDefault())
		assert.Equal(t, statusErr.CauseError.Error(), err.CauseError.Error())
	})

	t.Run("OutletOmzetGetData-ERROR", func(t *testing.T) {
		errNative := errors.New("failed to get data")
		statusErr := &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errNative,
		}
		gomock.InOrder(
			trxRepo.EXPECT().OutletOmzetGetTotal(&p).Return(int64(1), nil),
			trxRepo.EXPECT().OutletOmzetGetData(&p).Return(nil, statusErr),
		)
		res, err := trxService.OutletOmzetGet(&p)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.MessageOrDefault(), err.MessageOrDefault())
		assert.Equal(t, statusErr.CauseError.Error(), err.CauseError.Error())
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		expRes := make([]trxEntity.OutletOmzet, 0, 1)
		expRes = append(expRes, trxEntity.OutletOmzet{
			MerchantID:   1,
			MerchantName: "Merchant Name",
			OutletID:     1,
			OutletName:   "Outlet Name",
			Omzet:        500.75,
			Period:       "2024-03-07",
		})
		resInt64 := int64(len(expRes))
		gomock.InOrder(
			trxRepo.EXPECT().OutletOmzetGetTotal(&p).Return(resInt64, nil),
			trxRepo.EXPECT().OutletOmzetGetData(&p).Return(expRes, nil),
		)
		res, err := trxService.OutletOmzetGet(&p)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, resInt64, res.Total)
		assert.Equal(t, int(resInt64), len(res.Data))
		assert.Equal(t, expRes, res.Data)
	})
}

func setup(ctrl *gomock.Controller) (trxRepo *trxRepoMock.MockITransaction, ctx *ctx.Ctx, trxService *Service) {
	trxRepo = trxRepoMock.NewMockITransaction(ctrl)
	config, ctx := env()
	trxService = New(trxRepo, config)
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
