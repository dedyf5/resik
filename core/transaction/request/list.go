// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"net/http"

	ctx "github.com/dedyf5/resik/ctx"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	"github.com/dedyf5/resik/internal/identity"
	"github.com/dedyf5/resik/pkg/goku"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func (t *TransactionListGet) ToParam(c *ctx.Ctx, resolver identity.IdentityResolver) (param *paramTrx.TransactionsGet, err *resPkg.Status) {
	branchPublicIDs, errRes := resolver.GetBranchPublicIDsByPermission(c.Context, c.UserClaims().UserID(), "transaction:read")
	if errRes != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errRes)
	}

	orderStr := "-transacted_at"
	if t.Order != nil {
		orderStr = t.GetOrder()
	}

	return &paramTrx.TransactionsGet{
		Ctx:             c,
		BranchPublicIDs: branchPublicIDs,
		Filter:          *goku.NewFilter(t.GetSearch(), t.GetPage(), t.GetLimit()),
		Orders:          goku.OrdersBuilder(orderStr),
	}, nil
}
