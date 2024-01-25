// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package transaction

import "time"

type MerchantOmzet struct {
	MerchantName string    `gorm:"column:merchant_name"`
	Omzet        float64   `gorm:"column:omzet"`
	Date         time.Time `gorm:"column:date"`
}
