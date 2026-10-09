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
	branchEntity "github.com/dedyf5/resik/entities/branch"
	configEntity "github.com/dedyf5/resik/entities/config"
	orgEntity "github.com/dedyf5/resik/entities/organization"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	trxParam "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	identity "github.com/dedyf5/resik/internal/identity/mock"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	repoMock "github.com/dedyf5/resik/repositories/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"
)

func TestTransactionInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, trxRepo, _, ctx, trxService, _ := setup(ctrl)

	transaction := &trxEntity.Transaction{}

	t.Run("TestTransactionInsert TransactionInsert ERROR-500", func(t *testing.T) {
		okExpected := false
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			trxRepo.EXPECT().TransactionInsert(ctx, transaction).Return(okExpected, statusErr),
		)
		res, err := trxService.TransactionInsert(ctx, transaction)
		assert.Equal(t, okExpected, res)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestTransactionInsert ALL-SUCCESS", func(t *testing.T) {
		okExpected := true
		gomock.InOrder(
			trxRepo.EXPECT().TransactionInsert(ctx, transaction).Return(okExpected, nil),
		)
		res, err := trxService.TransactionInsert(ctx, transaction)
		assert.Equal(t, okExpected, res)
		assert.Nil(t, err)
	})
}

func TestTransactionUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orgRepo, trxRepo, _, ctx, trxService, resolver := setup(ctrl)

	orgPublicID, err := uuidPkg.NewUUIDV7()
	if err != nil {
		t.Error(err)
	}

	branchPublicID, err := uuidPkg.NewUUIDV7()
	if err != nil {
		t.Error(err)
	}

	branch := &branchEntity.Branch{
		PublicID: branchPublicID,
		Organization: orgEntity.Organization{
			PublicID: orgPublicID,
		},
	}

	transaction := &trxEntity.Transaction{
		BranchPublicID:       branchPublicID,
		OrganizationPublicID: orgPublicID,
	}

	t.Run("TestTransactionUpdate HasAccessByPublicID ERROR-500", func(t *testing.T) {
		okExpected := false
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			resolver.EXPECT().
				HasAccessByPublicID(
					ctx.Context,
					ctx.UserClaims().UserID(),
					"transaction:update",
					branchEntity.TABLE_NAME,
					branchPublicID,
				).
				Return(okExpected, errors.New("error 500")),
		)
		ok, err := trxService.TransactionUpdate(ctx, transaction)
		assert.Equal(t, okExpected, ok)
		assert.Equal(t, statusErr.Code, err.Code)
	})

	t.Run("TestTransactionUpdate BranchGetByPublicID ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			resolver.EXPECT().
				HasAccessByPublicID(
					ctx.Context,
					ctx.UserClaims().UserID(),
					"transaction:update",
					branchEntity.TABLE_NAME,
					branchPublicID,
				).
				Return(true, nil),
			orgRepo.EXPECT().BranchGetByPublicID(ctx, &branchPublicID).Return(nil, statusErr),
		)
		ok, err := trxService.TransactionUpdate(ctx, transaction)
		assert.False(t, ok)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestTransactionUpdate ALL-SUCCESS", func(t *testing.T) {
		okExpected := true
		gomock.InOrder(
			resolver.EXPECT().
				HasAccessByPublicID(
					ctx.Context,
					ctx.UserClaims().UserID(),
					"transaction:update",
					branchEntity.TABLE_NAME,
					branchPublicID,
				).
				Return(okExpected, nil),
			orgRepo.EXPECT().BranchGetByPublicID(ctx, &branchPublicID).Return(branch, nil),
			trxRepo.EXPECT().TransactionUpdate(ctx, transaction).Return(okExpected, nil),
		)
		ok, err := trxService.TransactionUpdate(ctx, transaction)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})
}

func TestHasAccessBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, _, _, ctx, trxService, resolver := setup(ctrl)

	branchPublicID, err := uuidPkg.NewUUIDV7()
	if err != nil {
		t.Error(err)
	}

	t.Run("TestHasAccessBranch HasAccessByPublicID ERROR-500", func(t *testing.T) {
		okExpected := false
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			resolver.EXPECT().
				HasAccessByPublicID(
					ctx.Context,
					ctx.UserClaims().UserID(),
					"transaction:update",
					branchEntity.TABLE_NAME,
					branchPublicID,
				).
				Return(okExpected, errors.New("error 500")),
		)
		ok, err := trxService.HasAccessBranch(ctx, &branchPublicID)
		assert.Equal(t, okExpected, ok)
		assert.Equal(t, statusErr.Code, err.Code)
	})

	t.Run("TestHasAccessBranch HasAccessByPublicID ERROR-403", func(t *testing.T) {
		okExpected := false
		statusErr := &resPkg.Status{
			Code: http.StatusForbidden,
		}
		gomock.InOrder(
			resolver.EXPECT().
				HasAccessByPublicID(
					ctx.Context,
					ctx.UserClaims().UserID(),
					"transaction:update",
					branchEntity.TABLE_NAME,
					branchPublicID,
				).
				Return(okExpected, nil),
		)
		ok, err := trxService.HasAccessBranch(ctx, &branchPublicID)
		assert.Equal(t, okExpected, ok)
		assert.Equal(t, statusErr.Code, err.Code)
	})

	t.Run("TestHasAccessBranch ALL-SUCCESS", func(t *testing.T) {
		okExpected := true
		gomock.InOrder(
			resolver.EXPECT().
				HasAccessByPublicID(
					ctx.Context,
					ctx.UserClaims().UserID(),
					"transaction:update",
					branchEntity.TABLE_NAME,
					branchPublicID,
				).
				Return(okExpected, nil),
		)
		ok, err := trxService.HasAccessBranch(ctx, &branchPublicID)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})
}

func TestValidateBranchOrganization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orgRepo, _, _, ctx, trxService, _ := setup(ctrl)

	orgPublicID, err := uuidPkg.NewUUIDV7()
	if err != nil {
		t.Error(err)
	}

	branchPublicID, err := uuidPkg.NewUUIDV7()
	if err != nil {
		t.Error(err)
	}

	branch := &branchEntity.Branch{
		PublicID: branchPublicID,
		Organization: orgEntity.Organization{
			PublicID: orgPublicID,
		},
	}

	t.Run("TestValidateBranchOrganization BranchGetByPublicID ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			orgRepo.EXPECT().BranchGetByPublicID(ctx, &branchPublicID).Return(nil, statusErr),
		)
		ok, err := trxService.ValidateBranchOrganization(ctx, &branchPublicID, &orgPublicID)
		assert.False(t, ok)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestValidateBranchOrganization BranchGetByPublicID ERROR-422-1", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusUnprocessableEntity,
		}
		gomock.InOrder(
			orgRepo.EXPECT().BranchGetByPublicID(ctx, &branchPublicID).Return(nil, nil),
		)
		ok, err := trxService.ValidateBranchOrganization(ctx, &branchPublicID, &orgPublicID)
		assert.False(t, ok)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.Code, err.Code)
	})

	t.Run("TestValidateBranchOrganization BranchGetByPublicID ERROR-422-2", func(t *testing.T) {
		orgPublicIDParam, uuidErr := uuidPkg.NewUUIDV7()
		if uuidErr != nil {
			t.Error(uuidErr)
		}
		statusErr := &resPkg.Status{
			Code: http.StatusUnprocessableEntity,
		}
		gomock.InOrder(
			orgRepo.EXPECT().BranchGetByPublicID(ctx, &branchPublicID).Return(branch, nil),
		)
		ok, err := trxService.ValidateBranchOrganization(ctx, &branchPublicID, &orgPublicIDParam)
		assert.False(t, ok)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.Code, err.Code)
	})

	t.Run("TestValidateBranchOrganization ALL-SUCCESS", func(t *testing.T) {
		gomock.InOrder(
			orgRepo.EXPECT().BranchGetByPublicID(ctx, &branchPublicID).Return(branch, nil),
		)
		ok, err := trxService.ValidateBranchOrganization(ctx, &branchPublicID, &orgPublicID)
		assert.True(t, ok)
		assert.Nil(t, err)
	})
}

func TestTransactionGetByPublicID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, trxRepo, userRepo, ctx, trxService, _ := setup(ctrl)

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

	_, trxRepo, userRepo, ctx, trxService, _ := setup(ctrl)

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

	_, trxRepo, _, ctx, trxService, _ := setup(ctrl)

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

	_, trxRepo, _, ctx, trxService, _ := setup(ctrl)

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

func setup(ctrl *gomock.Controller) (orgRepo *repoMock.MockIOrganization, trxRepo *repoMock.MockITransaction, userRepo *repoMock.MockIUser, ctx *ctx.Ctx, trxService *Service, resolver *identity.MockIdentityResolver) {
	config, ctx := env()
	resolver = identity.NewMockIdentityResolver(ctrl)
	orgRepo = repoMock.NewMockIOrganization(ctrl)
	trxRepo = repoMock.NewMockITransaction(ctrl)
	userRepo = repoMock.NewMockIUser(ctrl)
	trxService = New(config, resolver, orgRepo, trxRepo, userRepo)
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
