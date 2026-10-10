// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package dto

import (
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	userEntity "github.com/dedyf5/resik/entities/user"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

type Transaction struct {
	trxEntity.Transaction
	Creator *userEntity.User
	Updater *userEntity.User
}

func TransactionFromEntity(trx trxEntity.Transaction, users map[uuidPkg.UUIDV7]userEntity.User) Transaction {
	var creator, updater *userEntity.User
	if user, ok := users[trx.CreatedByPublicID]; ok {
		creator = &user
	}
	if user, ok := users[trx.UpdatedByPublicID]; ok {
		updater = &user
	}

	return Transaction{
		Transaction: trx,
		Creator:     creator,
		Updater:     updater,
	}
}

type Transactions []Transaction

func TransactionsFromEntity(trxs trxEntity.Transactions, users map[uuidPkg.UUIDV7]userEntity.User) Transactions {
	n := len(trxs)
	if n == 0 {
		return Transactions{}
	}

	res := make(Transactions, n)

	for i, trx := range trxs {
		res[i] = TransactionFromEntity(trx, users)
	}

	return res
}

type TransactionsResult struct {
	Data  Transactions
	Total int64
}

var (
	TransactionsResultEmpty = TransactionsResult{
		Data: Transactions{},
	}
)

type PERMISSION_CODE string

const (
	PERMISSION_CODE_CREATE PERMISSION_CODE = "transaction:create"
	PERMISSION_CODE_READ   PERMISSION_CODE = "transaction:read"
	PERMISSION_CODE_UPDATE PERMISSION_CODE = "transaction:update"
	PERMISSION_CODE_DELETE PERMISSION_CODE = "transaction:delete"
)

func (u PERMISSION_CODE) String() string {
	return string(u)
}
