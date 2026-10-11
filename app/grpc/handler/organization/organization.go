// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package organization

import (
	organizationService "github.com/dedyf5/resik/core/organization"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/internal/identity"
	validatorUtil "github.com/dedyf5/resik/utils/validator"
)

type OrganizationHandler struct {
	log                 *logCtx.Log
	validator           *validatorUtil.Validate
	resolver            identity.IdentityResolver
	organizationService organizationService.IService
}

func New(log *logCtx.Log, validator *validatorUtil.Validate, resolver identity.IdentityResolver, organizationService organizationService.IService) *OrganizationHandler {
	return &OrganizationHandler{
		log:                 log,
		validator:           validator,
		resolver:            resolver,
		organizationService: organizationService,
	}
}

func (h *OrganizationHandler) mustEmbedUnimplementedOrganizationServiceServer() {}
