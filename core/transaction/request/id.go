// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"net/http"

	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	permisEntity "github.com/dedyf5/resik/entities/permission"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	"github.com/dedyf5/resik/internal/identity"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

func (t *Id) ToParam(ctx *ctx.Ctx, resolver identity.IdentityResolver, permissionCode permisEntity.PERMISSION_CODE) (param *paramTrx.TransactionPublicID, err *resPkg.Status) {
	publicID, uuidErr := uuidPkg.ParseUUIDV7(t.GetId())
	if uuidErr != nil {
		return nil, resPkg.NewStatusMessage(
			http.StatusBadRequest,
			term.InvalidID.Localize(ctx.Lang().Localizer),
			uuidErr,
		)
	}

	branchPublicIDs, errRes := resolver.GetBranchPublicIDsByPermission(ctx.Context, ctx.UserClaims().UserID(), permissionCode.String())
	if errRes != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, errRes)
	}

	return &paramTrx.TransactionPublicID{
		Ctx:             ctx,
		BranchPublicIDs: branchPublicIDs,
		PublicID:        &publicID,
	}, nil
}
