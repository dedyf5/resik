// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"errors"
	"net/http"
	"slices"

	resPkg "github.com/dedyf5/resik/pkg/response"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

type MerchantOutletID struct {
	MerchantID       uint64         `gorm:"column:merchant_id"`
	MerchantPublicID uuidPkg.UUIDV7 `gorm:"column:merchant_public_id"`
	OutletID         uint64         `gorm:"column:outlet_id"`
	OutletPublicID   uuidPkg.UUIDV7 `gorm:"column:outlet_public_id"`
}

type MerchantOutletIDs []MerchantOutletID

func (mo MerchantOutletIDs) UniqueIDs() (merchantIDs []uint64, merchantPublicIDs []uuidPkg.UUIDV7, outletIDs []uint64, outletPublicIDs []uuidPkg.UUIDV7, err *resPkg.Status) {
	length := len(mo)

	MIDs := make([]uint64, 0, length)
	MPIDs := make([]uuidPkg.UUIDV7, 0, length)
	OIDs := make([]uint64, 0, length)
	OPIDs := make([]uuidPkg.UUIDV7, 0, length)

	for _, v := range mo {
		if v.OutletID > 0 {
			if v.OutletPublicID.IsEmpty() {
				err = resPkg.NewStatusError(http.StatusInternalServerError, errors.New("outlet public ID is empty"))
				return
			}
			OIDs = append(OIDs, v.OutletID)
			OPIDs = append(OPIDs, v.OutletPublicID)
		}

		if !slices.Contains(MIDs, v.MerchantID) {
			if v.MerchantPublicID.IsEmpty() {
				err = resPkg.NewStatusError(http.StatusInternalServerError, errors.New("merchant public ID is empty"))
				return
			}
			MIDs = append(MIDs, v.MerchantID)
			MPIDs = append(MPIDs, v.MerchantPublicID)
		}
	}

	return MIDs, MPIDs, OIDs, OPIDs, nil
}
