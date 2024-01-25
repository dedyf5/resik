// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	"time"

	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	transactionEntity "github.com/dedyf5/resik/entities/transaction"
	userEntity "github.com/dedyf5/resik/entities/user"
)

func (r *TransactionRepo) MerchantOmzet(merchantID int64, date time.Time) (*transactionEntity.MerchantOmzet, error) {
	return nil, nil
}

func (r *TransactionRepo) OutletOmzet(outletID int64, date time.Time) (*transactionEntity.OutletOmzet, error) {
	return nil, nil
}

func (r *TransactionRepo) GetMerchantByID(merchantID int64) (*merchantEntity.Merchant, error) {
	return nil, nil
}

func (r *TransactionRepo) GetMerchantByIDAndUserID(merchantID int64, userID int64) (*merchantEntity.Merchant, error) {
	return nil, nil
}

func (r *TransactionRepo) GetOutletByID(outletID int64) (*outletEntity.Outlet, error) {
	return nil, nil
}

func (r *TransactionRepo) GetOutletByIDAndCreatedBy(outletID int64, createdBy int64) (*outletEntity.Outlet, error) {
	return nil, nil
}

func (r *TransactionRepo) GetUserByUserNameAndPassword(userName string, password string) (*userEntity.User, error) {
	return nil, nil
}
