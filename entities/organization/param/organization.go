// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package param

import (
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/pkg/goku"
)

type OrganizationsGet struct {
	Ctx             *ctx.Ctx
	OrganizationIDs []uint64
	Filter          goku.Filter
	Orders          []goku.Order
}
