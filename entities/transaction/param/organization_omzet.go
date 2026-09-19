// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package param

import (
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/entities/groupperiod"
	"github.com/dedyf5/resik/pkg/goku"
)

type OrganizationOmzetGet struct {
	Ctx            *ctx.Ctx
	OrganizationID uint64
	GroupPeriod    groupperiod.GroupPeriod
	Filter         goku.Filter
	Orders         []goku.Order
}
