// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/dedyf5/resik/config"
	"github.com/dedyf5/resik/ctx"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	configEntity "github.com/dedyf5/resik/entities/config"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	trxParam "github.com/dedyf5/resik/entities/transaction/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	trxRepoMock "github.com/dedyf5/resik/repositories/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestMerchantOmzetGet(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	trxRepo := trxRepoMock.NewMockITransaction(ctl)
	config, ctx := env()
	trx := New(trxRepo, config)

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
		var resUint64 uint64 = 0
		gomock.InOrder(
			trxRepo.EXPECT().MerchantOmzetGetTotal(&p).Return(resUint64, statusErr),
		)
		res, err := trx.MerchantOmzetGet(&p)
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
		var resUint64 uint64 = 1
		gomock.InOrder(
			trxRepo.EXPECT().MerchantOmzetGetTotal(&p).Return(resUint64, nil),
			trxRepo.EXPECT().MerchantOmzetGetData(&p).Return(nil, statusErr),
		)
		res, err := trx.MerchantOmzetGet(&p)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.MessageOrDefault(), err.MessageOrDefault())
		assert.Equal(t, statusErr.CauseError.Error(), err.CauseError.Error())
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		expRes := []trxEntity.MerchantOmzet{}
		expRes = append(expRes, trxEntity.MerchantOmzet{
			MerchantID:   1,
			MerchantName: "Merchant Name",
			Omzet:        500.75,
			Period:       "2024-03-07",
		})
		var resUint64 uint64 = 1
		gomock.InOrder(
			trxRepo.EXPECT().MerchantOmzetGetTotal(&p).Return(resUint64, nil),
			trxRepo.EXPECT().MerchantOmzetGetData(&p).Return(expRes, nil),
		)
		res, err := trx.MerchantOmzetGet(&p)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, resUint64, res.Total)
		assert.Equal(t, int(resUint64), len(res.Data))
		assert.Equal(t, expRes, res.Data)
	})
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
		Lang:    langCtx.NewLang(language.English, &language.English, ""),
	}
	return
}
