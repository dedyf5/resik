// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package param

import (
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/pkg/goku"
)

type MerchantListGet struct {
	Ctx         *ctx.Ctx
	MerchantIDs []uint64
	Filter      goku.Filter
	Orders      []goku.Order
}
