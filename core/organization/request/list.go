// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"net/http"

	"github.com/dedyf5/resik/ctx"
	paramOrganization "github.com/dedyf5/resik/entities/organization/param"
	"github.com/dedyf5/resik/internal/identity"
	"github.com/dedyf5/resik/pkg/goku"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func (m *OrganizationListGet) ToParam(c *ctx.Ctx, resolver identity.IdentityResolver) (result *paramOrganization.OrganizationsGet, err *resPkg.Status) {
	orderStr := "name"
	if m.Order != nil {
		orderStr = m.GetOrder()
	}

	organizationIDs, errRes := resolver.GetOrganizationIDsByPermission(c.Context, c.UserClaims().UserID(), "organization:read")
	if errRes != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errRes)
	}

	return &paramOrganization.OrganizationsGet{
		Ctx:             c,
		OrganizationIDs: organizationIDs,
		Filter:          *goku.NewFilter(m.GetSearch(), m.GetPage(), m.GetLimit()),
		Orders:          goku.OrdersBuilder(orderStr),
	}, nil
}
