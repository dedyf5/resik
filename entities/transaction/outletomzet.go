// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

type OutletOmzet struct {
	MerchantID   int64   `gorm:"column:merchant_id"`
	MerchantName string  `gorm:"column:merchant_name"`
	OutletID     int64   `gorm:"column:outlet_id"`
	OutletName   string  `gorm:"column:outlet_name"`
	Omzet        float64 `gorm:"column:omzet"`
	Period       string  `gorm:"column:period"`
}
