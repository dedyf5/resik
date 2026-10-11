// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import (
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

type BranchOmzet struct {
	OrganizationPublicID uuidPkg.UUIDV7 `gorm:"column:organization_public_id"`
	OrganizationName     string         `gorm:"column:organization_name"`
	BranchPublicID       uuidPkg.UUIDV7 `gorm:"column:branch_public_id"`
	BranchName           string         `gorm:"column:branch_name"`
	Omzet                float64        `gorm:"column:omzet"`
	Period               string         `gorm:"column:period"`
}
