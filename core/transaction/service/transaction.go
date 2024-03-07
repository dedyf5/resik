// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	trxDTO "github.com/dedyf5/resik/core/transaction/dto"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func (s *Service) MerchantOmzetGet(param *paramTrx.MerchantOmzetGet) (res *trxDTO.MerchantOmzet, status *resPkg.Status) {
	total, status := s.transactionRepo.MerchantOmzetGetTotal(param)
	if status != nil {
		return nil, status
	}
	data, status := s.transactionRepo.MerchantOmzetGetData(param)
	if status != nil {
		return nil, status
	}
	return &trxDTO.MerchantOmzet{
		Data:  data,
		Total: total,
	}, nil
}

func (s *Service) OutletOmzetGet(param *paramTrx.OutletOmzetGet) (res *trxDTO.OutletOmzet, status *resPkg.Status) {
	total, status := s.transactionRepo.OutletOmzetGetTotal(param)
	if status != nil {
		return nil, status
	}
	data, status := s.transactionRepo.OutletOmzetGetData(param)
	if status != nil {
		return nil, status
	}
	return &trxDTO.OutletOmzet{
		Data:  data,
		Total: total,
	}, nil
}
