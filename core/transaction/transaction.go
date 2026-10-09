// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	trxDTO "github.com/dedyf5/resik/core/transaction/dto"
	"github.com/dedyf5/resik/ctx"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

//go:generate mockgen -source transaction.go -package mock -destination ./mock/transaction.go
type IService interface {
	TransactionInsert(ctx *ctx.Ctx, transaction *trxEntity.Transaction) (ok bool, err *resPkg.Status)
	TransactionUpdate(ctx *ctx.Ctx, trx *trxEntity.Transaction) (ok bool, err *resPkg.Status)
	HasAccessBranch(ctx *ctx.Ctx, branchPublicID *uuidPkg.UUIDV7) (ok bool, err *resPkg.Status)
	ValidateBranchOrganization(ctx *ctx.Ctx, branchPublicID, orgPublicID *uuidPkg.UUIDV7) (ok bool, err *resPkg.Status)
	TransactionGetByPublicID(param *paramTrx.TransactionGet) (res *trxDTO.Transaction, err *resPkg.Status)
	TransactionsGet(param *paramTrx.TransactionsGet) (res *trxDTO.TransactionsResult, err *resPkg.Status)
	OrganizationOmzetGet(param *paramTrx.OrganizationOmzetGet) (res *trxDTO.OrganizationOmzet, err *resPkg.Status)
	BranchOmzetGet(param *paramTrx.BranchOmzetGet) (res *trxDTO.BranchOmzet, err *resPkg.Status)
}
