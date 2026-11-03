// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/drivers"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/repositories/transaction/mysql"
	"github.com/dedyf5/resik/repositories/transaction/postgres"
	"gorm.io/gorm"
)

//go:generate mockgen -source transaction.go -package mock -destination ./mock/transaction.go
type ITransaction interface {
	TransactionInsert(ctx *ctx.Ctx, transaction *trxEntity.Transaction) (ok bool, err *resPkg.Status)
	TransactionUpdate(ctx *ctx.Ctx, transaction *trxEntity.Transaction) (ok bool, err *resPkg.Status)
	TransactionGetByPublicID(param *paramTrx.TransactionPublicID) (transaction *trxEntity.Transaction, err *resPkg.Status)
	TransactionsGetData(param *paramTrx.TransactionsGet) (transactions trxEntity.Transactions, err *resPkg.Status)
	TransactionsGetTotal(param *paramTrx.TransactionsGet) (total int64, err *resPkg.Status)
	TransactionDelete(param *paramTrx.TransactionPublicID) (ok bool, err *resPkg.Status)

	OrganizationOmzetGetData(param *paramTrx.OrganizationOmzetGet) (res []trxEntity.OrganizationOmzet, err *resPkg.Status)
	OrganizationOmzetGetTotal(param *paramTrx.OrganizationOmzetGet) (total int64, err *resPkg.Status)
	BranchOmzetGetData(param *paramTrx.BranchOmzetGet) (res []trxEntity.BranchOmzet, err *resPkg.Status)
	BranchOmzetGetTotal(param *paramTrx.BranchOmzetGet) (total int64, err *resPkg.Status)
}

func New(db *gorm.DB, engine drivers.SQLEngine) ITransaction {
	switch engine {
	case drivers.MySQL:
		return mysql.New(db)
	case drivers.PostgreSQL:
		return postgres.New(db)
	}
	return nil
}
