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
	OrganizationOmzetGet(param *paramTrx.OrganizationOmzetGet) (res *trxDTO.OrganizationOmzet, err *resPkg.Status)
	BranchOmzetGet(param *paramTrx.BranchOmzetGet) (res *trxDTO.BranchOmzet, err *resPkg.Status)
}
