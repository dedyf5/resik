// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	trxDTO "github.com/dedyf5/resik/core/transaction/dto"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

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
