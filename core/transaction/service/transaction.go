// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	trxDTO "github.com/dedyf5/resik/core/transaction/dto"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	resPkg "github.com/dedyf5/resik/pkg/response"
)

func (s *Service) MerchantOmzetGet(param *paramTrx.MerchantOmzetGet) (res *trxDTO.MerchantOmzet, err *resPkg.Status) {
	total, err := s.transactionRepo.MerchantOmzetGetTotal(param)
	if err != nil {
		return nil, err
	}
	data, err := s.transactionRepo.MerchantOmzetGetData(param)
	if err != nil {
		return nil, err
	}
	return &trxDTO.MerchantOmzet{
		Data:  data,
		Total: total,
	}, nil
}

func (s *Service) OutletOmzetGet(param *paramTrx.OutletOmzetGet) (res *trxDTO.OutletOmzet, err *resPkg.Status) {
	total, err := s.transactionRepo.OutletOmzetGetTotal(param)
	if err != nil {
		return nil, err
	}
	data, err := s.transactionRepo.OutletOmzetGetData(param)
	if err != nil {
		return nil, err
	}
	return &trxDTO.OutletOmzet{
		Data:  data,
		Total: total,
	}, nil
}
