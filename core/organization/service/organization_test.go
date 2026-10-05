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
	dtoOrganization "github.com/dedyf5/resik/core/organization/dto"
	"github.com/dedyf5/resik/ctx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	langCtx "github.com/dedyf5/resik/ctx/lang"
	"github.com/dedyf5/resik/ctx/log"
	configEntity "github.com/dedyf5/resik/entities/config"
	organizationEntity "github.com/dedyf5/resik/entities/organization"
	"github.com/dedyf5/resik/entities/organization/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	identityMock "github.com/dedyf5/resik/internal/identity/mock"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	repoMock "github.com/dedyf5/resik/repositories/mock"
)

func TestOrganizationInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resolver, ctx, organizationRepo, _, organizationService := setup(ctrl)

	userID := ctx.UserClaims().UserID()
	organization := &organizationEntity.Organization{}

	t.Run("TestOrganizationInsert OrganizationInsert ERROR-500", func(t *testing.T) {
		okExpected := false
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationInsert(ctx, organization).Return(okExpected, statusErr),
		)
		ok, err := organizationService.OrganizationInsert(ctx, organization)
		assert.Equal(t, okExpected, ok)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestOrganizationInsert InvalidateUserAccessOrganization ERROR-500", func(t *testing.T) {
		okExpected := true
		errExpected := errors.New("ERROR")
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationInsert(ctx, organization).Return(okExpected, nil),
			resolver.EXPECT().InvalidateUserAccessOrganization(ctx.Context, userID).Return(errExpected),
		)
		ok, err := organizationService.OrganizationInsert(ctx, organization)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})

	t.Run("TestOrganizationInsert ALL-SUCCESS", func(t *testing.T) {
		okExpected := true
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationInsert(ctx, organization).Return(okExpected, nil),
			resolver.EXPECT().InvalidateUserAccessOrganization(ctx.Context, userID).Return(nil),
		)
		ok, err := organizationService.OrganizationInsert(ctx, organization)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})
}

func TestOrganizationUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, ctx, organizationRepo, _, organizationService := setup(ctrl)

	organization := &organizationEntity.Organization{}

	t.Run("TestOrganizationUpdate OrganizationUpdate ERROR-500", func(t *testing.T) {
		okExpected := false
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationUpdate(ctx, organization).Return(okExpected, statusErr),
		)
		ok, err := organizationService.OrganizationUpdate(ctx, organization)
		assert.Equal(t, okExpected, ok)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestOrganizationUpdate ALL-SUCCESS", func(t *testing.T) {
		okExpected := true
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationUpdate(ctx, organization).Return(true, nil),
		)
		ok, err := organizationService.OrganizationUpdate(ctx, organization)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})
}

func TestOrganizationDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resolver, ctx, organizationRepo, _, organizationService := setup(ctrl)

	userID := ctx.UserClaims().UserID()
	organization := &organizationEntity.Organization{}

	t.Run("TestOrganizationDelete OrganizationDelete ERROR-500", func(t *testing.T) {
		okExpected := false
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationDelete(ctx, organization).Return(okExpected, statusErr),
		)
		ok, err := organizationService.OrganizationDelete(ctx, organization)
		assert.Equal(t, okExpected, ok)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestOrganizationDelete InvalidateUserAccessOrganization ERROR-500", func(t *testing.T) {
		okExpected := true
		errExpected := errors.New("ERROR")
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationDelete(ctx, organization).Return(okExpected, nil),
			resolver.EXPECT().InvalidateUserAccessOrganization(ctx.Context, userID).Return(errExpected),
		)
		ok, err := organizationService.OrganizationDelete(ctx, organization)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})

	t.Run("TestOrganizationDelete ALL-SUCCESS", func(t *testing.T) {
		okExpected := true
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationDelete(ctx, organization).Return(okExpected, nil),
			resolver.EXPECT().InvalidateUserAccessOrganization(ctx.Context, userID).Return(nil),
		)
		ok, err := organizationService.OrganizationDelete(ctx, organization)
		assert.Equal(t, okExpected, ok)
		assert.Nil(t, err)
	})
}

func TestOrganizationGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, ctx, organizationRepo, userRepo, organizationService := setup(ctrl)

	organization := &organizationEntity.Organization{}
	user := &userEntity.User{}
	users := userEntity.Users{
		*user,
	}
	userPublicIDs := []uuidPkg.UUIDV7{user.PublicID}

	t.Run("TestOrganizationGetByID OrganizationGetByID ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationGetByID(ctx, organization.ID).Return(organization, statusErr),
		)
		_, err := organizationService.OrganizationGetByID(ctx, organization.ID)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestOrganizationGetByID UsersGetByPublicIDs ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationGetByID(ctx, organization.ID).Return(organization, nil),
			userRepo.EXPECT().UsersGetByPublicIDs(ctx, userPublicIDs).Return(users, statusErr),
		)
		_, err := organizationService.OrganizationGetByID(ctx, organization.ID)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestOrganizationGetByID UsersGetByPublicIDs ERROR-404", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusNotFound,
		}
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationGetByID(ctx, organization.ID).Return(organization, nil),
			userRepo.EXPECT().UsersGetByPublicIDs(ctx, userPublicIDs).Return(nil, nil),
		)
		_, err := organizationService.OrganizationGetByID(ctx, organization.ID)
		assert.Equal(t, statusErr.Code, err.Code)
	})

	t.Run("TestOrganizationGetByID ALL-EMPTY", func(t *testing.T) {
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationGetByID(ctx, organization.ID).Return(nil, nil),
		)
		res, err := organizationService.OrganizationGetByID(ctx, organization.ID)
		assert.Nil(t, err)
		assert.Nil(t, res)
	})

	t.Run("TestOrganizationGetByID ALL-SUCCESS", func(t *testing.T) {
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationGetByID(ctx, organization.ID).Return(organization, nil),
			userRepo.EXPECT().UsersGetByPublicIDs(ctx, userPublicIDs).Return(users, nil),
		)
		res, err := organizationService.OrganizationGetByID(ctx, organization.ID)
		assert.Nil(t, err)
		assert.Equal(t, organization.PublicID, res.PublicID)
	})
}

func TestOrganizationsGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, ctx, organizationRepo, userRepo, organizationService := setup(ctrl)

	param := &param.OrganizationsGet{
		Ctx: ctx,
	}

	organizations := organizationEntity.Organizations{
		{
			ID:   1,
			Name: "organization1",
		},
	}

	user := &userEntity.User{}
	users := userEntity.Users{
		*user,
	}
	userPublicIDs := []uuidPkg.UUIDV7{user.PublicID}

	t.Run("TestOrganizationsGet OrganizationsGetTotal ERROR-500", func(t *testing.T) {
		var totalExpected int64 = 0
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationsGetTotal(param).Return(totalExpected, statusErr),
		)
		res, err := organizationService.OrganizationsGet(param)
		assert.Nil(t, res)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestOrganizationsGet OrganizationsGetTotal 0", func(t *testing.T) {
		var totalExpected int64 = 0
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationsGetTotal(param).Return(totalExpected, nil),
		)
		res, err := organizationService.OrganizationsGet(param)
		assert.Nil(t, err)
		assert.Equal(t, dtoOrganization.OrganizationsResultEmpty, *res)
	})

	var totalExpected int64 = 1
	t.Run("TestOrganizationsGet OrganizationsGetData ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationsGetTotal(param).Return(totalExpected, nil),
			organizationRepo.EXPECT().OrganizationsGetData(param).Return(nil, statusErr),
		)
		res, err := organizationService.OrganizationsGet(param)
		assert.Nil(t, res)
		assert.Equal(t, statusErr, err)
	})

	t.Run("TestOrganizationsGet OrganizationsGetData 0", func(t *testing.T) {
		organizationsEmpty := organizationEntity.Organizations{}
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationsGetTotal(param).Return(totalExpected, nil),
			organizationRepo.EXPECT().OrganizationsGetData(param).Return(organizationsEmpty, nil),
		)
		res, err := organizationService.OrganizationsGet(param)
		assert.Nil(t, err)
		assert.Equal(t, dtoOrganization.OrganizationsResultEmpty, *res)
	})

	t.Run("TestOrganizationsGet UsersGetByPublicIDs ERROR-500", func(t *testing.T) {
		statusErr := &resPkg.Status{
			Code: http.StatusInternalServerError,
		}
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationsGetTotal(param).Return(totalExpected, nil),
			organizationRepo.EXPECT().OrganizationsGetData(param).Return(organizations, nil),
			userRepo.EXPECT().UsersGetByPublicIDs(ctx, userPublicIDs).Return(nil, statusErr),
		)
		res, err := organizationService.OrganizationsGet(param)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, statusErr.Code, err.Code)
	})

	t.Run("TestOrganizationsGet ALL-SUCCESS", func(t *testing.T) {
		gomock.InOrder(
			organizationRepo.EXPECT().OrganizationsGetTotal(param).Return(totalExpected, nil),
			organizationRepo.EXPECT().OrganizationsGetData(param).Return(organizations, nil),
			userRepo.EXPECT().UsersGetByPublicIDs(ctx, userPublicIDs).Return(users, nil),
		)
		res, err := organizationService.OrganizationsGet(param)
		assert.Nil(t, err)
		assert.Len(t, res.Data, len(organizations))
		assert.Equal(t, organizations[0].ID, res.Data[0].ID)
		assert.Equal(t, organizations[0].Name, res.Data[0].Name)
	})
}

func setup(ctrl *gomock.Controller) (resolver *identityMock.MockIdentityResolver, ctx *ctx.Ctx, organizationRepo *repoMock.MockIOrganization, userRepo *repoMock.MockIUser, organizationService *Service) {
	resolver = identityMock.NewMockIdentityResolver(ctrl)
	organizationRepo = repoMock.NewMockIOrganization(ctrl)
	userRepo = repoMock.NewMockIUser(ctrl)
	config, ctx := env()
	organizationService = New(config, resolver, userRepo, organizationRepo)
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
