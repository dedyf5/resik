// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

type TransactionUpsert struct {
	Name      string `json:"name" validate:"required" example:"Shiina"`
	Address   string `json:"address" validate:"required" example:"Jl. Nusantara"`
	QTY       int64  `json:"qty" validate:"required,min=1" example:"3"`
	CreatedAt string `json:"created_at" validate:"required,datetime=2006-01-02 15:04:05" example:"2024-02-01 13:45:00"`
}
