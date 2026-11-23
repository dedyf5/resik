// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package param defines query and filter parameter structs for transaction and omzet operations.
package param

import (
	"time"

	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/pkg/goku"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

// TransactionPublicID holds execution context, allowed branch public IDs, and target transaction public ID
// for identifying a single transaction.
type TransactionPublicID struct {
	Ctx             *ctx.Ctx
	BranchPublicIDs []uuidPkg.UUIDV7
	PublicID        *uuidPkg.UUIDV7
}

// TransactionsGet holds execution context, branch public ID filters, datetime range boundaries,
// query filters, and ordering rules for retrieving transactions.
type TransactionsGet struct {
	Ctx             *ctx.Ctx
	BranchPublicIDs []uuidPkg.UUIDV7
	DatetimeStart   *time.Time
	DatetimeEnd     *time.Time
	Filter          goku.Filter
	Orders          []goku.Order
}
