// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

type MerchantOmzetGet struct {
	ID    uint64 `param:"id" validate:"required,min=1" example:"1"`
	Page  *uint  `query:"page" validate:"omitempty,min=1" example:"1"`
	Limit *uint  `query:"limit" validate:"omitempty,min=1,max=100" example:"10"`
}

type TransactionUpsert struct {
	Name      string `json:"name" validate:"required" example:"Shiina"`
	Address   string `json:"address" validate:"required" example:"Jl. Nusantara"`
	QTY       int64  `json:"qty" validate:"required,min=1" example:"3"`
	CreatedAt string `json:"created_at" validate:"required,datetime=2006-01-02 15:04:05" example:"2024-02-01 13:45:00"`
}

func (s *MerchantOmzetGet) PageOrDefault() uint {
	if s.Page != nil {
		if *s.Page > 0 {
			return *s.Page
		}
	}
	return PageDefault
}

func (s *MerchantOmzetGet) LimitOrDefault() uint {
	if s.Limit != nil {
		if *s.Limit > 0 {
			return *s.Limit
		}
	}
	return LimitDefault
}
