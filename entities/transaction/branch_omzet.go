// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

type BranchOmzet struct {
	OrganizationID   uint64  `gorm:"column:organization_id"`
	OrganizationName string  `gorm:"column:organization_name"`
	BranchID         uint64  `gorm:"column:branch_id"`
	BranchName       string  `gorm:"column:branch_name"`
	Omzet            float64 `gorm:"column:omzet"`
	Period           string  `gorm:"column:period"`
}
