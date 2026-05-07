// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"slices"

	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

type MerchantOutletID struct {
	MerchantID       uint64         `gorm:"column:merchant_id"`
	MerchantPublicID uuidPkg.UUIDV7 `gorm:"column:merchant_public_id"`
	OutletID         uint64         `gorm:"column:outlet_id"`
	OutletPublicID   uuidPkg.UUIDV7 `gorm:"column:outlet_public_id"`
}

type MerchantOutletIDs []MerchantOutletID

func (mo MerchantOutletIDs) UniqueIDs() (merchantIDs []uint64, outletIDs []uint64) {
	length := len(mo)

	MIDs := make([]uint64, 0, length)
	OIDs := make([]uint64, 0, length)

	for _, v := range mo {
		if v.OutletID > 0 {
			OIDs = append(OIDs, v.OutletID)
		}

		if !slices.Contains(MIDs, v.MerchantID) {
			MIDs = append(MIDs, v.MerchantID)
		}
	}

	return MIDs, OIDs
}
