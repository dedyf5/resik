// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"net/http"
	"time"

	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"github.com/dedyf5/resik/utils/datetime"
)

func (p *TransactionPost) ToEntity(ctx *ctx.Ctx) (res *trxEntity.Transaction, err *resPkg.Status) {
	transactionAt, err := datetime.FromString(p.GetTransactedAt(), time.RFC3339, ctx)
	if err != nil {
		return nil, err
	}

	organizationPublicID, uuidErr := uuidPkg.ParseUUIDV7(p.GetBranch().GetOrganization().GetId())
	if uuidErr != nil {
		return nil, invalidID(ctx, uuidErr)
	}

	branch := p.GetBranch()
	org := branch.GetOrganization()

	branchPublicID, uuidErr := uuidPkg.ParseUUIDV7(branch.GetId())
	if uuidErr != nil {
		return nil, invalidID(ctx, uuidErr)
	}

	userPublicID := ctx.UserClaims().UserPublicID()
	now := time.Now()

	return &trxEntity.Transaction{
		OrganizationPublicID: organizationPublicID,
		OrganizationName:     org.GetName(),
		BranchPublicID:       branchPublicID,
		BranchName:           branch.GetName(),
		BillTotal:            p.GetBillTotal(),
		TransactedAt:         *transactionAt,
		CreatedAt:            now,
		CreatedByPublicID:    userPublicID,
		UpdatedAt:            now,
		UpdatedByPublicID:    userPublicID,
	}, nil
}

func invalidID(ctx *ctx.Ctx, err error) *resPkg.Status {
	return resPkg.NewStatusMessage(
		http.StatusBadRequest,
		term.InvalidID.Localize(ctx.Lang().Localizer),
		err,
	)
}
