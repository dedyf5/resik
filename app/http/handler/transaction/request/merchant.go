// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package request

type GeMerchantOmzet struct {
	ID    uint64 `param:"id" validate:"required,min=1" example:"1"`
	Page  *uint  `query:"page" validate:"omitempty,min=1" example:"1"`
	Limit *uint  `query:"limit" validate:"omitempty,min=1,max=100" example:"10"`
}

func (s *GeMerchantOmzet) PageOrDefault() uint {
	if s.Page != nil {
		if *s.Page > 0 {
			return *s.Page
		}
	}
	return PageDefault
}

func (s *GeMerchantOmzet) LimitOrDefault() uint {
	if s.Limit != nil {
		if *s.Limit > 0 {
			return *s.Limit
		}
	}
	return LimitDefault
}
