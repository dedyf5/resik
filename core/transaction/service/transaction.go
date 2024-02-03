// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package service

import (
	"errors"
	"fmt"
	"time"

	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	trxEntity "github.com/dedyf5/resik/entities/transaction"
	paramTrx "github.com/dedyf5/resik/entities/transaction/param"
	userEntity "github.com/dedyf5/resik/entities/user"
	statusPkg "github.com/dedyf5/resik/pkg/status"
)

const (
	perPage = 10
)

func (s *Service) MerchantOmzetGet(param *paramTrx.MerchantOmzetGet) (res []trxEntity.MerchantOmzet, status *statusPkg.Status) {
	return s.transactionRepo.MerchantOmzetGet(param)
}

func (s *Service) OutletOmzet(outletID int64, dates []time.Time) ([]trxEntity.OutletOmzet, error) {
	var outlets []trxEntity.OutletOmzet
	base, err1 := s.transactionRepo.GetOutletByID(outletID)
	if err1 != nil {
		return outlets, err1
	}
	for _, v := range dates {
		result, err2 := s.transactionRepo.OutletOmzet(outletID, v)
		if err2 == nil {
			outlets = append(outlets, *result)
		} else {
			outlet := trxEntity.OutletOmzet{
				MerchantName: base.Merchant.MerchantName,
				OutletName:   base.OutletName,
				Omzet:        0,
				Date:         v,
			}
			outlets = append(outlets, outlet)
		}
	}
	return outlets, nil
}

func (s *Service) GetUserByUserNameAndPassword(userName string, password string) (*userEntity.User, error) {
	return s.transactionRepo.GetUserByUserNameAndPassword(userName, password)
}

func (s *Service) ValidateAuthRequest(username, password string) error {
	if username == "" {
		err := errors.New("the user name is empty")
		return err
	}
	if password == "" {
		err := errors.New("the password is empty")
		return err
	}
	return nil
}

func (s *Service) ValidateMerchantUser(merchantID int64, userID int64) (*merchantEntity.Merchant, error) {
	return s.transactionRepo.GetMerchantByIDAndUserID(merchantID, userID)
}

func (s *Service) ValidateOutletUser(outletID int64, createdBy int64) (*outletEntity.Outlet, error) {
	return s.transactionRepo.GetOutletByIDAndCreatedBy(outletID, createdBy)
}

func (s *Service) Dates(date *time.Time, page int) []time.Time {
	var dates []time.Time
	last := date.AddDate(0, 1, -1)
	total := last.Day()
	var start int
	if page == 1 {
		start = 0
	} else {
		start = (page - 1) * perPage
	}
	end := start + perPage
	if end > total {
		end = total
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	if start < total {
		for day := start; day < end; day++ {
			newDate := date.AddDate(0, 0, day).In(loc)
			hour := newDate.Hour()
			if hour > 0 {
				min := hour * -1
				f := fmt.Sprint(min) + `h`
				minDur, _ := time.ParseDuration(f)
				newDate = newDate.Add(minDur)
			}
			dates = append(dates, newDate)
		}
	}
	return dates
}
