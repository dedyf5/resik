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
	"time"

	"github.com/dedyf5/resik/config"
	dtoTrx "github.com/dedyf5/resik/core/transaction/dto"
	"github.com/dedyf5/resik/ctx"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/log"
	configEntity "github.com/dedyf5/resik/entities/config"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	trxParam "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	repoMock "github.com/dedyf5/resik/repositories/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"
)

func TestTransactionGetByPublicID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trxRepo, userRepo, ctx, trxService := setup(ctrl)

	param := &trxParam.TransactionGet{
		Ctx: ctx,
	}

	transaction := &trxEntity.Transaction{}
	user := &userEntity.User{}
	users := userEntity.Users{
		*user,
	}
	userPublicIDs := []uuidPkg.UUIDV7{user.PublicID}

	t.Run("TestTransactionGetByPublicID TransactionGetByPublicID ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			trxRepo.EXPECT().TransactionGetByPublicID(param).Return(transaction, statusErr),
		)
		_, err := trxService.TransactionGetByPublicID(param)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestTransactionGetByPublicID UsersGetByPublicIDs ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			trxRepo.EXPECT().TransactionGetByPublicID(param).Return(transaction, nil),
			userRepo.EXPECT().UsersGetByPublicIDs(ctx, userPublicIDs).Return(users, statusErr),
		)
		_, err := trxService.TransactionGetByPublicID(param)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestTransactionGetByPublicID UsersGetByPublicIDs ERROR-404", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusNotFound,
		}
		gomock.InOrder(
			trxRepo.EXPECT().TransactionGetByPublicID(param).Return(transaction, nil),
			userRepo.EXPECT().UsersGetByPublicIDs(ctx, userPublicIDs).Return(nil, nil),
		)
		_, err := trxService.TransactionGetByPublicID(param)
		assert.Equal(t, statusErr.Code, err.Code)
	})

	t.Run("TestTransactionGetByPublicID ALL-EMPTY", func(t *testing.T) {
		gomock.InOrder(
			trxRepo.EXPECT().TransactionGetByPublicID(param).Return(nil, nil),
		)
		res, err := trxService.TransactionGetByPublicID(param)
		assert.Nil(t, err)
		assert.Nil(t, res)
	})

	t.Run("TestTransactionGetByPublicID ALL-SUCCESS", func(t *testing.T) {
		gomock.InOrder(
			trxRepo.EXPECT().TransactionGetByPublicID(param).Return(transaction, nil),
			userRepo.EXPECT().UsersGetByPublicIDs(ctx, userPublicIDs).Return(users, nil),
		)
		res, err := trxService.TransactionGetByPublicID(param)
		assert.Nil(t, err)
		assert.Equal(t, transaction.PublicID, res.PublicID)
	})
}

func TestTransactionsGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trxRepo, userRepo, ctx, trxService := setup(ctrl)

	param := &trxParam.TransactionsGet{
		Ctx: ctx,
	}

	datetime := time.Now()

	transactions := trxEntity.Transactions{
		{
			ID:           1,
			TransactedAt: datetime,
			UpdatedAt:    datetime,
			CreatedAt:    datetime,
		},
	}

	user := &userEntity.User{}
	users := userEntity.Users{
		*user,
	}
	userPublicIDs := []uuidPkg.UUIDV7{user.PublicID}

	t.Run("TestTransactionsGet TransactionsGetTotal ERROR-500", func(t *testing.T) {
		var totalExpected int64 = 0
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			trxRepo.EXPECT().TransactionsGetTotal(param).Return(totalExpected, statusErr),
		)
		res, err := trxService.TransactionsGet(param)
		assert.Nil(t, res)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestTransactionsGet TransactionsGetTotal 0", func(t *testing.T) {
		var totalExpected int64 = 0
		gomock.InOrder(
			trxRepo.EXPECT().TransactionsGetTotal(param).Return(totalExpected, nil),
		)
		res, err := trxService.TransactionsGet(param)
		assert.Nil(t, err)
		assert.Equal(t, dtoTrx.TransactionsResultEmpty, *res)
	})

	var totalExpected int64 = 1
	t.Run("TestTransactionsGet TransactionsGetData ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			trxRepo.EXPECT().TransactionsGetTotal(param).Return(totalExpected, nil),
			trxRepo.EXPECT().TransactionsGetData(param).Return(nil, statusErr),
		)
		res, err := trxService.TransactionsGet(param)
		assert.Nil(t, res)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestTransactionsGet TransactionsGetData 0", func(t *testing.T) {
		transactionsEmpty := trxEntity.Transactions{}
		gomock.InOrder(
			trxRepo.EXPECT().TransactionsGetTotal(param).Return(totalExpected, nil),
			trxRepo.EXPECT().TransactionsGetData(param).Return(transactionsEmpty, nil),
		)
		res, err := trxService.TransactionsGet(param)
		assert.Nil(t, err)
		assert.Equal(t, dtoTrx.TransactionsResultEmpty, *res)
	})

	t.Run("TestTransactionsGet UsersGetByPublicIDs ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			trxRepo.EXPECT().TransactionsGetTotal(param).Return(totalExpected, nil),
			trxRepo.EXPECT().TransactionsGetData(param).Return(transactions, nil),
			userRepo.EXPECT().UsersGetByPublicIDs(ctx, userPublicIDs).Return(nil, statusErr),
		)
		res, err := trxService.TransactionsGet(param)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.Code, err.Code)
	})

	t.Run("TestTransactionsGet ALL-SUCCESS", func(t *testing.T) {
		gomock.InOrder(
			trxRepo.EXPECT().TransactionsGetTotal(param).Return(totalExpected, nil),
			trxRepo.EXPECT().TransactionsGetData(param).Return(transactions, nil),
			userRepo.EXPECT().UsersGetByPublicIDs(ctx, userPublicIDs).Return(users, nil),
		)
		res, err := trxService.TransactionsGet(param)
		assert.Nil(t, err)
		assert.Len(t, res.Data, len(transactions))
		assert.Equal(t, transactions[0].ID, res.Data[0].Transaction.ID)
		assert.Equal(t, transactions[0].BillTotal, res.Data[0].BillTotal)
		assert.Equal(t, transactions[0].TransactedAt, res.Data[0].TransactedAt)
	})
}

func TestOrganizationOmzetGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trxRepo, _, ctx, trxService := setup(ctrl)

	organizationPublicID, err := uuidPkg.NewUUIDV7()
	if err != nil {
		t.Fatal(err)
	}

	p := trxParam.OrganizationOmzetGet{
		Ctx:                  ctx,
		OrganizationPublicID: organizationPublicID,
	}

	t.Run("TestOrganizationOmzetGet OrganizationOmzetGetTotal ERROR-500", func(t *testing.T) {
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

	t.Run("TestOrganizationOmzetGet OrganizationOmzetGetData ERROR-500", func(t *testing.T) {
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

	t.Run("TestOrganizationOmzetGet ALL-SUCCESS", func(t *testing.T) {
		expRes := make([]trxEntity.OrganizationOmzet, 0, 1)
		expRes = append(expRes, trxEntity.OrganizationOmzet{
			OrganizationPublicID: organizationPublicID,
			OrganizationName:     "Organization Name",
			Omzet:                500.75,
			Period:               "2024-03-07",
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

	trxRepo, _, ctx, trxService := setup(ctrl)

	organizationPublicID, err := uuidPkg.NewUUIDV7()
	if err != nil {
		t.Fatal(err)
	}

	branchPublicID, err := uuidPkg.NewUUIDV7()
	if err != nil {
		t.Fatal(err)
	}

	p := trxParam.BranchOmzetGet{
		Ctx:            ctx,
		BranchPublicID: branchPublicID,
	}

	t.Run("TestBranchOmzetGet BranchOmzetGetTotal ERROR-500", func(t *testing.T) {
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

	t.Run("TestBranchOmzetGet BranchOmzetGetData ERROR-500", func(t *testing.T) {
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

	t.Run("TestBranchOmzetGet ALL-SUCCESS", func(t *testing.T) {
		expRes := make([]trxEntity.BranchOmzet, 0, 1)
		expRes = append(expRes, trxEntity.BranchOmzet{
			OrganizationPublicID: organizationPublicID,
			OrganizationName:     "Organization Name",
			BranchPublicID:       branchPublicID,
			BranchName:           "Branch Name",
			Omzet:                500.75,
			Period:               "2024-03-07",
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

func setup(ctrl *gomock.Controller) (trxRepo *repoMock.MockITransaction, userRepo *repoMock.MockIUser, ctx *ctx.Ctx, trxService *Service) {
	trxRepo = repoMock.NewMockITransaction(ctrl)
	userRepo = repoMock.NewMockIUser(ctrl)
	config, ctx := env()
	trxService = New(config, trxRepo, userRepo)
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
