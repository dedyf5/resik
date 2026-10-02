// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package param

import (
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/pkg/goku"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

type TransactionsGet struct {
	Ctx             *ctx.Ctx
	BranchPublicIDs []uuidPkg.UUIDV7
	Filter          goku.Filter
	Orders          []goku.Order
}
