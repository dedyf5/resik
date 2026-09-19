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

func TestOrganizationOmzetGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trxRepo, ctx, trxService := setup(ctrl)

	p := trxParam.OrganizationOmzetGet{
		Ctx:            ctx,
		OrganizationID: 1,
	}

	t.Run("OrganizationOmzetGetTotal-ERROR", func(t *testing.T) {
		errNative := errors.New("failed to get total")
		statusErr := &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errNative,
		}
		gomock.InOrder(
			trxRepo.EXPECT().OrganizationOmzetGetTotal(&p).Return(int64(0), statusErr),
		)
		res, err := trxService.OrganizationOmzetGet(&p)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.MessageOrDefault(), err.MessageOrDefault())
		assert.Equal(t, statusErr.CauseError.Error(), err.CauseError.Error())
	})

	t.Run("OrganizationOmzetGetData-ERROR", func(t *testing.T) {
		errNative := errors.New("failed to get data")
		statusErr := &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errNative,
		}
		gomock.InOrder(
			trxRepo.EXPECT().OrganizationOmzetGetTotal(&p).Return(int64(1), nil),
			trxRepo.EXPECT().OrganizationOmzetGetData(&p).Return(nil, statusErr),
		)
		res, err := trxService.OrganizationOmzetGet(&p)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.MessageOrDefault(), err.MessageOrDefault())
		assert.Equal(t, statusErr.CauseError.Error(), err.CauseError.Error())
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		expRes := make([]trxEntity.OrganizationOmzet, 0, 1)
		expRes = append(expRes, trxEntity.OrganizationOmzet{
			OrganizationID:   1,
			OrganizationName: "Organization Name",
			Omzet:            500.75,
			Period:           "2024-03-07",
		})
		resInt64 := int64(len(expRes))
		gomock.InOrder(
			trxRepo.EXPECT().OrganizationOmzetGetTotal(&p).Return(resInt64, nil),
			trxRepo.EXPECT().OrganizationOmzetGetData(&p).Return(expRes, nil),
		)
		res, err := trxService.OrganizationOmzetGet(&p)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, resInt64, res.Total)
		assert.Len(t, res.Data, int(resInt64))
		assert.Equal(t, expRes, res.Data)
	})
}

func TestBranchOmzetGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trxRepo, ctx, trxService := setup(ctrl)

	p := trxParam.BranchOmzetGet{
		Ctx:      ctx,
		BranchID: 1,
	}

	t.Run("BranchOmzetGetTotal-ERROR", func(t *testing.T) {
		errNative := errors.New("failed to get total")
		statusErr := &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errNative,
		}
		gomock.InOrder(
			trxRepo.EXPECT().BranchOmzetGetTotal(&p).Return(int64(0), statusErr),
		)
		res, err := trxService.BranchOmzetGet(&p)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.MessageOrDefault(), err.MessageOrDefault())
		assert.Equal(t, statusErr.CauseError.Error(), err.CauseError.Error())
	})

	t.Run("BranchOmzetGetData-ERROR", func(t *testing.T) {
		errNative := errors.New("failed to get data")
		statusErr := &resPkg.Status{
			Code:       http.StatusInternalServerError,
			CauseError: errNative,
		}
		gomock.InOrder(
			trxRepo.EXPECT().BranchOmzetGetTotal(&p).Return(int64(1), nil),
			trxRepo.EXPECT().BranchOmzetGetData(&p).Return(nil, statusErr),
		)
		res, err := trxService.BranchOmzetGet(&p)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.MessageOrDefault(), err.MessageOrDefault())
		assert.Equal(t, statusErr.CauseError.Error(), err.CauseError.Error())
	})

	t.Run("ALL-SUCCESS", func(t *testing.T) {
		expRes := make([]trxEntity.BranchOmzet, 0, 1)
		expRes = append(expRes, trxEntity.BranchOmzet{
			OrganizationID:   1,
			OrganizationName: "Organization Name",
			BranchID:         1,
			BranchName:       "Branch Name",
			Omzet:            500.75,
			Period:           "2024-03-07",
		})
		resInt64 := int64(len(expRes))
		gomock.InOrder(
			trxRepo.EXPECT().BranchOmzetGetTotal(&p).Return(resInt64, nil),
			trxRepo.EXPECT().BranchOmzetGetData(&p).Return(expRes, nil),
		)
		res, err := trxService.BranchOmzetGet(&p)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, resInt64, res.Total)
		assert.Len(t, res.Data, int(resInt64))
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
		Module: configEntity.Module{
			Name:        "REST",
			NameKey:     "rest",
			LangDefault: language.English,
			Host:        "0.0.0.0",
			Port:        8081,
		},
	}
	context := context.WithValue(context.Background(), langCtx.ContextKey, langCtx.NewLangLocaleDir(language.English, &language.English, "", fmt.Sprintf("%s%s", "../../../", langCtx.LocaleDir)))
	c, _ = ctx.NewCtx(context, &log.Log{})
	return
}
