// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	trxDTO "github.com/dedyf5/resik/core/transaction/dto"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

//go:generate mockgen -source transaction.go -package mock -destination ./mock/transaction.go
type IService interface {
	MerchantOmzetGet(param *paramTrx.MerchantOmzetGet) (res *trxDTO.MerchantOmzet, status *resPkg.Status)
	OutletOmzetGet(param *paramTrx.OutletOmzetGet) (res *trxDTO.OutletOmzet, status *resPkg.Status)
}
