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
	permisEntity "github.com/dedyf5/resik/entities/permission"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

func (s *Service) TransactionInsert(ctx *ctx.Ctx, trx *trxEntity.Transaction) (ok bool, err *resPkg.Status) {
	_, err = s.HasAccessBranch(ctx, &trx.BranchPublicID, permisEntity.PERMISSION_CODE_TRX_CREATE)
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
	permissionCode := permisEntity.PERMISSION_CODE_TRX_UPDATE

	branchPublicIDs, err := s.resolver.GetBranchPublicIDsByPermission(
		ctx.Context,
		ctx.UserClaims().UserID(),
		permissionCode.String(),
	)
	if err != nil {
		return false, err
	}
	if len(branchPublicIDs) == 0 {
		return false, resPkg.NewStatusMessage(
			http.StatusForbidden,
			term.AccessDenied.Localize(ctx.Lang().Localizer),
			nil,
		)
	}

	trxBefore, err := s.transactionRepo.TransactionGetByPublicID(&paramTrx.TransactionPublicID{
		Ctx:             ctx,
		BranchPublicIDs: branchPublicIDs,
		PublicID:        &trx.PublicID,
	})
	if err != nil {
		return false, err
	}
	if trxBefore == nil {
		return false, resPkg.NewStatusMessage(
			http.StatusNotFound,
			term.NotFoundVal.Localize(
				ctx.Lang().Localizer,
				term.Transaction.Localize(ctx.Lang().Localizer),
			),
			nil,
		)
	}

	_, err = s.HasAccessBranch(ctx, &trx.BranchPublicID, permissionCode)
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

func (s *Service) HasAccessBranch(ctx *ctx.Ctx, branchPublicID *uuidPkg.UUIDV7, permissionCode permisEntity.PERMISSION_CODE) (ok bool, err *resPkg.Status) {
	return s.resolver.HasAccessByPublicID(
		ctx.Context, ctx.Lang(), ctx.UserClaims().UserID(),
		permissionCode.String(),
		branch.TABLE_NAME,
		*branchPublicID,
	)
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
