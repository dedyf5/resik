// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"net/http"

	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	"github.com/dedyf5/resik/internal/identity"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

func (t *TransactionDetailGet) ToParam(ctx *ctx.Ctx, resolver identity.IdentityResolver) (param *paramTrx.TransactionGet, err *resPkg.Status) {
	publicID, uuidErr := uuidPkg.ParseUUIDV7(t.GetId())
	if uuidErr != nil {
		return nil, resPkg.NewStatusMessage(
			http.StatusBadRequest,
			term.InvalidID.Localize(ctx.Lang().Localizer),
			uuidErr,
		)
	}

	branchPublicIDs, errRes := resolver.GetBranchPublicIDsByPermission(ctx.Context, ctx.UserClaims().UserID(), "transaction:read")
	if errRes != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errRes)
	}

	return &paramTrx.TransactionGet{
		Ctx:             ctx,
		BranchPublicIDs: branchPublicIDs,
		PublicID:        publicID,
	}, nil
}
