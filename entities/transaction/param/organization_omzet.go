// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package param

import (
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/groupperiod"
	"github.com/dedyf5/resik/pkg/goku"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

// OrganizationOmzetGet holds execution context, organization public ID, grouping period, filters,
// and ordering rules for retrieving organization omzet data.
type OrganizationOmzetGet struct {
	Ctx                  *ctx.Ctx
	OrganizationPublicID uuidPkg.UUIDV7
	GroupPeriod          groupperiod.GroupPeriod
	Filter               goku.Filter
	Orders               []goku.Order
}
