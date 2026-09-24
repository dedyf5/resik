// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"net/http"

	dtoOrganization "github.com/dedyf5/resik/core/organization/dto"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	organizationEntity "github.com/dedyf5/resik/entities/organization"
	paramOrganization "github.com/dedyf5/resik/entities/organization/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func (s *Service) OrganizationInsert(ctx *ctx.Ctx, organization *organizationEntity.Organization) (ok bool, err *resPkg.Status) {
	ok, err = s.organizationRepo.OrganizationInsert(ctx, organization)
	if err != nil {
		return ok, err
	}

	userID := ctx.UserClaims().UserID()

	if err := s.resolver.InvalidateUserAccessOrganization(ctx.Context, userID); err != nil {
		ctx.Log().Errorf("failed to invalidate user access organization for user id: %d and error: %s", userID, err.Error())
	}

	return ok, nil
}

func (s *Service) OrganizationUpdate(ctx *ctx.Ctx, organization *organizationEntity.Organization) (ok bool, err *resPkg.Status) {
	return s.organizationRepo.OrganizationUpdate(ctx, organization)
}

func (s *Service) OrganizationGetByID(ctx *ctx.Ctx, organizationID uint64) (*dtoOrganization.Organization, *resPkg.Status) {
	organization, err := s.organizationRepo.OrganizationGetByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization == nil {
		return nil, nil
	}

	users, err := s.userRepo.UsersGetByPublicIDs(ctx, organization.UniqueAllUserPublicIDs())
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		localizer := ctx.Lang().Localizer
		return nil, resPkg.NewStatusMessage(
			http.StatusNotFound,
			term.NotFoundVal.Localize(localizer, term.User.Localize(localizer)),
			nil,
		)
	}

	res := dtoOrganization.OrganizationFromEntity(*organization, users.UniquePublicIDsMap())
	return &res, nil
}

func (s *Service) OrganizationsGet(param *paramOrganization.OrganizationsGet) (res *dtoOrganization.OrganizationsResult, err *resPkg.Status) {
	total, err := s.organizationRepo.OrganizationsGetTotal(param)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &dtoOrganization.OrganizationsResultEmpty, nil
	}

	organizations, err := s.organizationRepo.OrganizationsGetData(param)
	if err != nil {
		return nil, err
	}

	if len(organizations) == 0 {
		return &dtoOrganization.OrganizationsResultEmpty, nil
	}

	users, err := s.userRepo.UsersGetByPublicIDs(param.Ctx, organizations.UniqueAllUserPublicIDs())
	if err != nil {
		return nil, err
	}

	return &dtoOrganization.OrganizationsResult{
		Data:  dtoOrganization.OrganizationsFromEntity(organizations, users.UniquePublicIDsMap()),
		Total: total,
	}, nil
}

func (s *Service) OrganizationDelete(ctx *ctx.Ctx, organization *organizationEntity.Organization) (ok bool, err *resPkg.Status) {
	ok, err = s.organizationRepo.OrganizationDelete(ctx, organization)
	if err != nil {
		return ok, err
	}

	userID := ctx.UserClaims().UserID()

	if err := s.resolver.InvalidateUserAccessOrganization(ctx.Context, userID); err != nil {
		ctx.Log().Errorf("failed to invalidate user access organization for user id: %d and error: %s", userID, err.Error())
	}

	return ok, nil
}
