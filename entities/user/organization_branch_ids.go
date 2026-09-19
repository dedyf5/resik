// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package user

import (
	"slices"

	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

type OrganizationBranchID struct {
	OrganizationID       uint64         `gorm:"column:organization_id"`
	OrganizationPublicID uuidPkg.UUIDV7 `gorm:"column:organization_public_id"`
	BranchID             uint64         `gorm:"column:branch_id"`
	BranchPublicID       uuidPkg.UUIDV7 `gorm:"column:branch_public_id"`
}

type OrganizationBranchIDs []OrganizationBranchID

func (mo OrganizationBranchIDs) UniqueIDs() (organizationIDs []uint64, branchIDs []uint64) {
	length := len(mo)

	MIDs := make([]uint64, 0, length)
	OIDs := make([]uint64, 0, length)

	for _, v := range mo {
		if v.BranchID > 0 {
			OIDs = append(OIDs, v.BranchID)
		}

		if !slices.Contains(MIDs, v.OrganizationID) {
			MIDs = append(MIDs, v.OrganizationID)
		}
	}

	return MIDs, OIDs
}
