// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package dto

import trxEntity "github.com/dedyf5/resik/entities/transaction"

type BranchOmzet struct {
	Data  []trxEntity.BranchOmzet
	Total int64
}
