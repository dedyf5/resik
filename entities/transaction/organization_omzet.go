// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

type OrganizationOmzet struct {
	OrganizationID   uint64  `gorm:"column:organization_id"`
	OrganizationName string  `gorm:"column:organization_name"`
	Omzet            float64 `gorm:"column:omzet"`
	Period           string  `gorm:"column:period"`
}
