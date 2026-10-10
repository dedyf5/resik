// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"net/http"

	trxDTO "github.com/dedyf5/resik/core/transaction/dto"
	"github.com/dedyf5/resik/ctx"
	"github.com/dedyf5/resik/ctx/lang/term"
	"github.com/dedyf5/resik/entities/branch"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

func (s *Service) TransactionInsert(ctx *ctx.Ctx, trx *trxEntity.Transaction) (ok bool, err *resPkg.Status) {
	_, err = s.HasAccessBranch(ctx, &trx.BranchPublicID, trxDTO.PERMISSION_CODE_CREATE)
	if err != nil {
		return false, err
	}

	_, err = s.ValidateBranchOrganization(ctx, &trx.BranchPublicID, &trx.OrganizationPublicID)
	if err != nil {
		return false, err
	}

	ok, err = s.transactionRepo.TransactionInsert(ctx, trx)
	return
}

func (s *Service) TransactionUpdate(ctx *ctx.Ctx, trx *trxEntity.Transaction) (ok bool, err *resPkg.Status) {
	_, err = s.HasAccessBranch(ctx, &trx.BranchPublicID, trxDTO.PERMISSION_CODE_UPDATE)
	if err != nil {
		return false, err
	}

	_, err = s.ValidateBranchOrganization(ctx, &trx.BranchPublicID, &trx.OrganizationPublicID)
	if err != nil {
		return false, err
	}

	ok, err = s.transactionRepo.TransactionUpdate(ctx, trx)
	return
}

func (s *Service) HasAccessBranch(ctx *ctx.Ctx, branchPublicID *uuidPkg.UUIDV7, permissionCode trxDTO.PERMISSION_CODE) (ok bool, err *resPkg.Status) {
	hasAccess, accessErr := s.resolver.HasAccessByPublicID(
		ctx.Context, ctx.UserClaims().UserID(),
		permissionCode.String(),
		branch.TABLE_NAME,
		*branchPublicID,
	)
	if accessErr != nil {
		return false, resPkg.NewStatusError(http.StatusInternalServerError, accessErr)
	}
	if !hasAccess {
		return false, resPkg.NewStatusError(
			http.StatusForbidden,
			accessErr,
		)
	}
	return true, nil
}

func (s *Service) ValidateBranchOrganization(ctx *ctx.Ctx, branchPublicID, orgPublicID *uuidPkg.UUIDV7) (ok bool, err *resPkg.Status) {
	branch, err := s.organizationRepo.BranchGetByPublicID(ctx, branchPublicID)
	if err != nil {
		return false, err
	}
	if branch == nil {
		return false, resPkg.NewStatusError(
			http.StatusUnprocessableEntity,
			nil,
		)
	}
	return branch.IsBelongToOrganization(orgPublicID)
}

func (s *Service) TransactionGetByPublicID(param *paramTrx.TransactionPublicID) (res *trxDTO.Transaction, err *resPkg.Status) {
	transaction, err := s.transactionRepo.TransactionGetByPublicID(param)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, nil
	}

	users, err := s.userRepo.UsersGetByPublicIDs(param.Ctx, transaction.UniqueAllUserPublicIDs())
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		localizer := param.Ctx.Lang().Localizer
		return nil, resPkg.NewStatusMessage(
			http.StatusNotFound,
			term.NotFoundVal.Localize(localizer, term.User.Localize(localizer)),
			nil,
		)
	}

	resDTO := trxDTO.TransactionFromEntity(*transaction, users.UniquePublicIDsMap())
	return &resDTO, nil
}

func (s *Service) TransactionsGet(param *paramTrx.TransactionsGet) (res *trxDTO.TransactionsResult, err *resPkg.Status) {
	total, err := s.transactionRepo.TransactionsGetTotal(param)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &trxDTO.TransactionsResultEmpty, nil
	}

	transactions, err := s.transactionRepo.TransactionsGetData(param)
	if err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return &trxDTO.TransactionsResultEmpty, nil
	}

	users, err := s.userRepo.UsersGetByPublicIDs(param.Ctx, transactions.UniqueAllUserPublicIDs())
	if err != nil {
		return nil, err
	}

	return &trxDTO.TransactionsResult{
		Data:  trxDTO.TransactionsFromEntity(transactions, users.UniquePublicIDsMap()),
		Total: total,
	}, nil
}

func (s *Service) TransactionDelete(param *paramTrx.TransactionPublicID) (ok bool, err *resPkg.Status) {
	trx, err := s.transactionRepo.TransactionGetByPublicID(param)
	if err != nil {
		return false, err
	}
	if trx == nil {
		return false, resPkg.NewStatusError(
			http.StatusNotFound,
			nil,
		)
	}
	return s.transactionRepo.TransactionDelete(param)
}

func (s *Service) OrganizationOmzetGet(param *paramTrx.OrganizationOmzetGet) (res *trxDTO.OrganizationOmzet, err *resPkg.Status) {
	total, err := s.transactionRepo.OrganizationOmzetGetTotal(param)
	if err != nil {
		return nil, err
	}
	data, err := s.transactionRepo.OrganizationOmzetGetData(param)
	if err != nil {
		return nil, err
	}
	return &trxDTO.OrganizationOmzet{
		Data:  data,
		Total: total,
	}, nil
}

func (s *Service) BranchOmzetGet(param *paramTrx.BranchOmzetGet) (res *trxDTO.BranchOmzet, err *resPkg.Status) {
	total, err := s.transactionRepo.BranchOmzetGetTotal(param)
	if err != nil {
		return nil, err
	}
	data, err := s.transactionRepo.BranchOmzetGetData(param)
	if err != nil {
		return nil, err
	}
	return &trxDTO.BranchOmzet{
		Data:  data,
		Total: total,
	}, nil
}
