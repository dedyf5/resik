// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

import (
	"net/http"
	"time"

	commonCore "github.com/dedyf5/resik/core/common"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"github.com/dedyf5/resik/utils/datetime"
)

type transactionUpsert interface {
	GetBranch() *commonCore.Branch
	GetBillTotal() float64
	GetTransactedAt() string
}

func toEntity[T transactionUpsert](p T, ctx *ctx.Ctx) (res *trxEntity.Transaction, err *resPkg.Status) {
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

	return &trxEntity.Transaction{
		OrganizationPublicID: organizationPublicID,
		OrganizationName:     org.GetName(),
		BranchPublicID:       branchPublicID,
		BranchName:           branch.GetName(),
		BillTotal:            p.GetBillTotal(),
		TransactedAt:         *transactionAt,
	}, nil
}

func (p *TransactionPost) ToEntity(ctx *ctx.Ctx) (res *trxEntity.Transaction, err *resPkg.Status) {
	res, err = toEntity(p, ctx)
	if err != nil {
		return nil, err
	}

	userPublicID := ctx.UserClaims().UserPublicID()
	now := time.Now()

	res.CreatedAt = now
	res.CreatedByPublicID = userPublicID
	res.UpdatedAt = now
	res.UpdatedByPublicID = userPublicID

	return
}

func (p *TransactionPut) ToEntity(ctx *ctx.Ctx) (res *trxEntity.Transaction, err *resPkg.Status) {
	res, err = toEntity(p, ctx)
	if err != nil {
		return nil, err
	}

	publicID, uuidErr := uuidPkg.ParseUUIDV7(p.GetId())
	if uuidErr != nil {
		return nil, invalidID(ctx, uuidErr)
	}

	userPublicID := ctx.UserClaims().UserPublicID()
	now := time.Now()

	res.PublicID = publicID
	res.UpdatedAt = now
	res.UpdatedByPublicID = userPublicID

	return
}

func invalidID(ctx *ctx.Ctx, err error) *resPkg.Status {
	return resPkg.NewStatusMessage(
		http.StatusBadRequest,
		term.InvalidID.Localize(ctx.Lang().Localizer),
		err,
	)
}
