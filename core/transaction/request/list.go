// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"time"

	ctx "github.com/dedyf5/resik/ctx"
	permisEntity "github.com/dedyf5/resik/entities/permission"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	"github.com/dedyf5/resik/internal/identity"
	"github.com/dedyf5/resik/pkg/goku"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/datetime"
)

func (t *TransactionListGet) ToParam(c *ctx.Ctx, resolver identity.IdentityResolver) (param *paramTrx.TransactionsGet, err *resPkg.Status) {
	branchPublicIDs, errRes := resolver.GetBranchPublicIDsByPermission(c.Context, c.UserClaims().UserID(), permisEntity.PERMISSION_CODE_TRX_READ.String())
	if errRes != nil {
		return nil, errRes
	}

	datetimeStart, errRes := datetime.FromString(t.GetDatetimeStart(), time.RFC3339, c)
	if errRes != nil {
		return nil, errRes
	}

	datetimeEnd, errRes := datetime.FromString(t.GetDatetimeEnd(), time.RFC3339, c)
	if errRes != nil {
		return nil, errRes
	}

	orderStr := "-transacted_at"
	if t.Order != nil {
		orderStr = t.GetOrder()
	}

	return &paramTrx.TransactionsGet{
		Ctx:             c,
		BranchPublicIDs: branchPublicIDs,
		DatetimeStart:   datetimeStart,
		DatetimeEnd:     datetimeEnd,
		Filter:          *goku.NewFilter(t.GetSearch(), t.GetPage(), t.GetLimit()),
		Orders:          goku.OrdersBuilder(orderStr),
	}, nil
}
