// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package dto

import trxEntity "github.com/dedyf5/resik/entities/transaction"

type OutletOmzet struct {
	Data  []trxEntity.OutletOmzet
	Total uint64
}
