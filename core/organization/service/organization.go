// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"net/http"

	dtoOrganization "github.com/dedyf5/resik/core/organization/dto"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	paramOrganization "github.com/dedyf5/resik/entities/organization/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

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
