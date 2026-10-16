// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package param

import (
	"time"

	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/pkg/goku"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

type TransactionPublicID struct {
	Ctx             *ctx.Ctx
	BranchPublicIDs []uuidPkg.UUIDV7
	PublicID        *uuidPkg.UUIDV7
}

type TransactionsGet struct {
	Ctx             *ctx.Ctx
	BranchPublicIDs []uuidPkg.UUIDV7
	DatetimeStart   *time.Time
	DatetimeEnd     *time.Time
	Filter          goku.Filter
	Orders          []goku.Order
}
